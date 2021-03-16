/*
 * Copyright 2019-2021 VMware, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package api

import (
	"errors"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// Chart API struct
type Chart struct {
}

// Router is Chart router definition method
func (c *Chart) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	chart := r.Group("/chart")
	chart.Use(authMiddleware.MiddlewareFunc())
	{
		chart.POST("", c.createChart)
		chart.GET("/", c.getChartList)
		chart.GET("/:chartId", c.getChart)
		chart.DELETE("/:chartId", c.deleteChart)
	}
}

// HelmUUID HelmUUID
type HelmUUID struct {
	HelmUUID string
}

// createChart Upload a Chart
// @Summary Upload a Chart
// @Tags Chart
// @Produce  json
// @Param file formData file true "Cluster Chart"
// @Success 200 {object} JSONResult{data=HelmUUID} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /chart [post]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Chart) createChart(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("FormFile error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Debug().Str("Filename", file.Filename).Msg("upload file")

	f, err := file.Open()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	chartRequested, err := loader.LoadArchive(f)
	if err != nil {
		log.Error().Err(err).Msg("LoadArchive error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	hc := modules.HelmChart{}
	helmChart, err := hc.ChartRequestedToHelmChart(chartRequested)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = helmChart.Upload()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("helmChart.Chart", helmChart.Chart).Msg("createChart Success")
	c.JSON(200, gin.H{"msg": "CreateChart Success", "data": gin.H{"helmUUID": helmChart.Uuid}})
}

// getChartList List all historical Charts
// @Summary List all historical Charts
// @Tags Chart
// @Produce  json
// @Success 200 {object} JSONResult{data=[]modules.HelmChart} "Success"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /chart [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Chart) getChartList(c *gin.Context) {

	chartList, err := new(modules.HelmChart).GetList()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("chartList", chartList).Msg("getChartList Success")
	c.JSON(200, gin.H{"msg": "getChartList Success", "data": chartList})
}

// getChart Get Chart by chartId
// @Summary Get Chart by chartId
// @Tags Chart
// @Produce  json
// @Param  chartId path string true "Chart Id"
// @Success 200 {object} JSONResult{data=modules.HelmChart} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /chart/{chartId} [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Chart) getChart(c *gin.Context) {

	chartID := c.Param("chartId")
	if chartID == "" {
		log.Error().Err(errors.New("not exit chartId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit chartId"})
		return
	}

	hc := modules.HelmChart{Uuid: chartID}
	chartList, err := hc.Get()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("chartList", chartList).Msg("getChart Success")
	c.JSON(200, gin.H{"msg": "getChart Success", "data": chartList})
}

// deleteChart Delete Chart by chartId
// @Summary Delete Chart by chartId
// @Tags Chart
// @Produce  json
// @Param  chartId path string true "chart Id"
// @Success 200 {object} JSONEMSGResult "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONEMSGResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /chart/{chartId} [delete]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Chart) deleteChart(c *gin.Context) {
	chartID := c.Param("chartId")
	if chartID == "" {
		c.JSON(400, gin.H{"error": "not exit chartId"})
		return
	}
	chart := new(modules.HelmChart)
	_, err := chart.DeleteByUuid(chartID)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("result", "deleteChart Success").Msg("deleteChart Success")
	c.JSON(200, gin.H{"msg": "deleteChart Success"})
}
