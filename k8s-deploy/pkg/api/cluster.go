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

	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/job"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Cluster api of Cluster
type Cluster struct {
}

// Router is Cluster router definition method
func (c *Cluster) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	cluster := r.Group("/cluster")
	cluster.Use(authMiddleware.MiddlewareFunc())
	{
		cluster.POST("", c.createCluster)
		cluster.PUT("", c.setCluster)
		cluster.GET("/", c.getClusterList)
		cluster.GET("/:clusterId", c.getCluster)
		cluster.DELETE("/:clusterId", c.deleteCluster)
	}
}

// createCluster Create a new Cluster
// @Summary Create a new Cluster
// @Tags Cluster
// @Produce  json
// @Param  clusterArgs body modules.ClusterArgs true "Cluster Args"
// @Success 200 {object} JSONResult{data=modules.Job} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /cluster [post]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
func (*Cluster) createCluster(c *gin.Context) {

	user, _ := c.Get(identityKey)

	clusterArgs := new(modules.ClusterArgs)

	if err := c.ShouldBindJSON(&clusterArgs); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("parameters", clusterArgs).Msg("parameters")

	// create job and use goroutine do job result save to db
	j, err := job.ClusterInstall(clusterArgs, user.(*User).Username)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", j).Msg("createCluster Success")
	c.JSON(200, gin.H{"msg": "createCluster Success", "data": j})
}

// setCluster Updates a Cluster in the store with form data
// @Summary Updates a Cluster in the store with form data
// @Tags Cluster
// @Produce  json
// @Param  ClusterArgs body modules.ClusterArgs true "Cluster Args"
// @Success 200 {object} JSONResult{data=modules.Job} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /cluster [put]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
func (*Cluster) setCluster(c *gin.Context) {

	//cluster := new(db.Cluster)
	//if err := c.ShouldBindJSON(&cluster); err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}

	user, _ := c.Get(identityKey)

	clusterArgs := new(modules.ClusterArgs)

	if err := c.ShouldBindJSON(&clusterArgs); err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("parameters", clusterArgs).Msg("parameters")

	// create job and use goroutine do job result save to db
	j, err := job.ClusterUpdate(clusterArgs, user.(*User).Username)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", j).Msg("setCluster Success")
	c.JSON(200, gin.H{"msg": "setCluster Success", "data": j})
}

// getCluster create a Cluster
// @Summary create Cluster
// @Tags Cluster
// @Produce  json
// @Param   clusterId   path  string  true  "Cluster Id"
// @Success 200 {object} JSONResult{data=modules.Cluster} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /cluster/{clusterId} [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
func (*Cluster) getCluster(c *gin.Context) {

	clusterID := c.Param("clusterId")
	if clusterID == "" {
		log.Error().Err(errors.New("not exit clusterId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit clusterId"})
		return
	}

	hc := modules.Cluster{Uuid: clusterID}
	cluster, err := hc.Get()
	if err != nil {
		log.Error().Err(err).Str("uuid", clusterID).Msg("get Cluster error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	cluster.Info, err = service.GetClusterInfo(cluster.Name, cluster.NameSpace)
	if err != nil {
		log.Error().Err(err).Str("Name", cluster.Name).Str("NameSpace", cluster.NameSpace).Msg("GetClusterInfo error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if cluster.Spec == nil {
		cluster.Spec = make(map[string]interface{})
	}

	// Check the cluster status and update to DB.
	// If cluster's Status isn't Running or Unavailable, don't check the cluster status.
	if cluster.Status == modules.ClusterStatusRunning || cluster.Status ==
		modules.ClusterStatusUnavailable {

		if service.CheckClusterInfoStatus(cluster.Info) {
			if cluster.Status != modules.ClusterStatusRunning {
				dbErr := cluster.SetStatus(modules.ClusterStatusRunning)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job setStatus error")
				}
			}

		} else {
			if cluster.Status != modules.ClusterStatusUnavailable {
				dbErr := cluster.SetStatus(modules.ClusterStatusUnavailable)
				if dbErr != nil {
					log.Error().Err(dbErr).Msg("job setStatus error")
				}
			}
		}
	}

	log.Debug().Interface("data", cluster).Msg("getCluster Success")
	c.JSON(200, gin.H{"msg": "getCluster Success", "data": &cluster})
}

// getClusterList List all available Clusters
// @Summary List all available Clusters
// @Tags Cluster
// @Produce  json
// @Param   all     query    boolean      true        "get All Cluster" default(false)
// @Success 200 {object} JSONResult{data=[]modules.Cluster} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /cluster/ [get]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
func (*Cluster) getClusterList(c *gin.Context) {

	all := false

	if c.Query("all") == "true" {
		all = true
	}

	log.Debug().Bool("all", all).Msg("get args")

	clusterList, err := new(modules.Cluster).GetListAll(all)

	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", clusterList).Msg("getClusterList Success")
	c.JSON(200, gin.H{"msg": "getClusterList Success", "data": clusterList})
}

// deleteCluster Delete a Cluster
// @Summary Delete a Cluster
// @Tags Cluster
// @Produce  json
// @Param   clusterId   path  string  true  "Cluster Id"
// @Success 200 {object} JSONResult{data=modules.Job} "Success"
// @Failure 400 {object} JSONERRORResult "Bad Request"
// @Failure 401 {object} JSONERRORResult "Unauthorized operation"
// @Failure 500 {object} JSONERRORResult "Internal server error"
// @Router /cluster/{clusterId} [delete]
// @Param Authorization header string true "Authentication header"
// @Security ApiKeyAuth
// @Security OAuth2Application[write, admin]
func (*Cluster) deleteCluster(c *gin.Context) {

	user, _ := c.Get(identityKey)

	clusterID := c.Param("clusterId")
	if clusterID == "" {
		log.Error().Err(errors.New("not exit clusterId")).Msg("request error")
		c.JSON(400, gin.H{"error": "not exit clusterId"})
	}

	j, err := job.ClusterDelete(clusterID, user.(*User).Username)
	if err != nil {
		log.Error().Err(err).Msg("request error")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("data", j).Msg("deleteCluster Success")
	c.JSON(200, gin.H{"msg": "deleteCluster Success", "data": j})
}
