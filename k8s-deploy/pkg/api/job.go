/*
 * Copyright 2019-2020 VMware, Inc.
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
)

// Job API struct
type Job struct {
}

// Router of Job API
func (j *Job) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	job := r.Group("/job")
	job.Use(authMiddleware.MiddlewareFunc())
	{
		job.GET("/", j.getJobList)
		job.GET("/:jobId", j.getJob)
		job.DELETE("/:jobId", j.deleteJob)
	}
}

// getJobList List all historical jobs
// @Summary List all historical jobs
// @Tags Job
// @Produce  json
// @Success 200 {object} JSONResult{data=[]modules.Job} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Router /job/ [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Job) getJobList(c *gin.Context) {

	jobList, err := new(modules.Job).GetList()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", jobList).Msg("getJobList success")
	c.JSON(200, gin.H{"data": jobList, "msg": "getJobList success"})
}

// getJob Get job by jobId
// @Summary Get job by jobId
// @Tags Job
// @Produce  json
// @Param  jobId path int true "Job ID"
// @Success 200 {object} JSONResult{data=[]modules.Job} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /job/{jobId} [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Job) getJob(c *gin.Context) {

	jobID := c.Param("jobId")
	if jobID == "" {
		log.Error().Err(errors.New("not exit jobId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit jobId"})
		return
	}
	j := modules.Job{Uuid: jobID}
	job, err := j.Get()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", job).Msg("getJob success")
	c.JSON(200, gin.H{"msg": "getJob success", "data": job})
}

// deleteJob Delete Job by jobId
// @Summary Delete Job by jobId
// @Tags Job
// @Produce  json
// @Param  jobId path int true "Job ID"
// @Success 200 {object} JSONEMSGResult "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /job/{jobId} [delete]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*Job) deleteJob(c *gin.Context) {

	jobID := c.Param("jobId")
	if jobID == "" {
		log.Error().Err(errors.New("not exit jobId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit jobId"})
		return
	}
	j := modules.Job{Uuid: jobID}
	_, err := j.Delete()
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("msg", "delete Job success").Msg("delete Job success")
	c.JSON(200, gin.H{"msg": "delete Job success"})
}
