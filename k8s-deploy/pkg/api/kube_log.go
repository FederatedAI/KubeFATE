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
	"time"

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type kubeLog struct {
}

// RequestArgs Request Args
type RequestArgs struct {
	Container                    string    `form:"container"`
	Follow                       bool      `form:"follow"`
	Previous                     bool      `form:"previous"`
	SinceSeconds                 *int64    `form:"since"`
	SinceTime                    time.Time `form:"since-time" time_format:"2006-01-02T15:04:05Z07:00"`
	Timestamps                   bool      `form:"timestamps"`
	TailLines                    *int64    `form:"tail"`
	LimitBytes                   *int64    `form:"limit-bytes"`
	InsecureSkipTLSVerifyBackend bool
}

func (e *kubeLog) Router(r *gin.RouterGroup) {
	authMiddleware, _ := GetAuthMiddleware()
	kubeLog := r.Group("/log")
	kubeLog.Use(authMiddleware.MiddlewareFunc())
	{
		kubeLog.GET("/:clusterID", e.getClusterLog)
		kubeLog.GET("/:clusterID/ws", e.getClusterLogWs)
	}
}

// getClusterLog Get Cluster log by clusterID
// @Summary Get Cluster log by clusterID
// @Description When the container of requestArgs is not set, the logs of all components will be obtained
// @Tags Log
// @Produce  json
// @Param  clusterID path int true "Cluster ID"
// @Param  container query string true "container name"
// @Param  previous query boolean  true "previous" default(false)
// @Param  since query string true "since"
// @Param  since-time query string true "since-time"
// @Param  timestamps query boolean  true "timestamps" default(false)
// @Param  tail query string true "tail"
// @Param  limit-bytes query string true "limit-bytes"
// @Success 200 {object} JSONResult{data=string} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /log/{clusterID} [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*kubeLog) getClusterLog(c *gin.Context) {

	clusterID := c.Param("clusterID")
	if clusterID == "" {
		log.Error().Err(errors.New("not exit clusterID")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit clusterID"})
		return
	}

	requestArgs := new(RequestArgs)
	if err := c.ShouldBind(&requestArgs); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hc := modules.Cluster{Uuid: clusterID}
	cluster, err := hc.Get()
	if err != nil {
		log.Error().Err(err).Str("uuid", clusterID).Msg("get Cluster error")
		c.JSON(400, gin.H{"error": "get Cluster error, " + err.Error()})
		return
	}

	buf, err := service.GetLogs(&service.LogChanArgs{
		Name:                         cluster.Name,
		Namespace:                    cluster.NameSpace,
		Container:                    requestArgs.Container,
		Follow:                       false,
		Previous:                     requestArgs.Previous,
		SinceSeconds:                 requestArgs.SinceSeconds,
		SinceTime:                    requestArgs.SinceTime,
		Timestamps:                   requestArgs.Timestamps,
		TailLines:                    requestArgs.TailLines,
		LimitBytes:                   requestArgs.LimitBytes,
		InsecureSkipTLSVerifyBackend: requestArgs.InsecureSkipTLSVerifyBackend,
	})

	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Debug().Int("data.size", buf.Len()).Msg("getClusterLog Success")
	c.JSON(200, gin.H{"data": buf.String(), "msg": "getClusterLog Success"})

}

// getClusterLog Get Cluster Log flow (http1.1)
// @Summary Get Cluster Log flow (http1.1)
// @Tags Log
// @Produce  json
// @Param  clusterID path int true "Cluster ID"
// @Param  container query string true "container name"
// @Param  previous query boolean  true "previous"
// @Param  since query string true "since"
// @Param  since-time query string true "since-time"
// @Param  timestamps query boolean  true "timestamps"
// @Param  tail query string true "tail"
// @Param  limit-bytes query string true "limit-bytes"
// @Success 200 {object} string "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /log/{clusterID}/ws [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
func (*kubeLog) getClusterLogWs(c *gin.Context) {

	clusterID := c.Param("clusterID")
	if clusterID == "" {
		log.Error().Err(errors.New("not exit clusterID")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit clusterID"})
		return
	}

	requestArgs := new(RequestArgs)
	if err := c.ShouldBind(&requestArgs); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hc := modules.Cluster{Uuid: clusterID}
	cluster, err := hc.Get()
	if err != nil {
		log.Error().Err(err).Str("uuid", clusterID).Msg("get Cluster error")
		c.JSON(400, gin.H{"error": "get Cluster error, " + err.Error()})
		return
	}

	handler := websocket.Handler(func(c *websocket.Conn) {
		log.Debug().Msg("get log websocket reader Success")
		defer log.Debug().Msg("websocket close")

		err := service.WriteLog(c, &service.LogChanArgs{
			Name:                         cluster.Name,
			Namespace:                    cluster.NameSpace,
			Container:                    requestArgs.Container,
			Follow:                       true,
			Previous:                     requestArgs.Previous,
			SinceSeconds:                 requestArgs.SinceSeconds,
			SinceTime:                    requestArgs.SinceTime,
			Timestamps:                   requestArgs.Timestamps,
			TailLines:                    requestArgs.TailLines,
			LimitBytes:                   requestArgs.LimitBytes,
			InsecureSkipTLSVerifyBackend: requestArgs.InsecureSkipTLSVerifyBackend,
		})
		log.Warn().Err(err).Msg("writeLog err, if the log stream is closed, you can ignore this prompt")
	})
	handler.ServeHTTP(c.Writer, c.Request)
}
