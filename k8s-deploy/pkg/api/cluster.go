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
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/job"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/modules"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Cluster struct {
}

// Router is cluster router definition method
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

func (_ *Cluster) createCluster(c *gin.Context) {

	user, _ := c.Get(identityKey)

	clusterArgs := new(job.ClusterArgs)

	if err := c.ShouldBindJSON(&clusterArgs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("parameters", clusterArgs).Msg("parameters")

	// create job and use goroutine do job result save to db
	j, err := job.ClusterInstall(clusterArgs, user.(*User).Username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "createCluster success", "data": j})
}

func (_ *Cluster) setCluster(c *gin.Context) {

	//cluster := new(db.Cluster)
	//if err := c.ShouldBindJSON(&cluster); err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}

	user, _ := c.Get(identityKey)

	clusterArgs := new(job.ClusterArgs)

	if err := c.ShouldBindJSON(&clusterArgs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Interface("parameters", clusterArgs).Msg("parameters")

	// create job and use goroutine do job result save to db
	j, err := job.ClusterUpdate(clusterArgs, user.(*User).Username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "setCluster success", "data": j})
}

func (_ *Cluster) getCluster(c *gin.Context) {

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		c.JSON(400, gin.H{"error": "not exit clusterId"})
		return
	}

	hc := modules.Cluster{Uuid: clusterId}
	cluster, err := hc.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	cluster.Info, err = service.GetClusterInfo(cluster.Name, cluster.NameSpace)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	if cluster.Spec == nil {
		cluster.Spec = make(map[string]interface{})
	}

	c.JSON(200, gin.H{"data": &cluster})
}

func (_ *Cluster) getClusterList(c *gin.Context) {

	all := false
	qall := c.Query("all")
	if qall == "true" {
		all = true
	}

	log.Debug().Bool("all", all).Msg("get args")

	clusterList, err := new(modules.Cluster).GetListAll(all)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "getClusterList success", "data": clusterList})
}

func (_ *Cluster) deleteCluster(c *gin.Context) {

	user, _ := c.Get(identityKey)

	clusterId := c.Param("clusterId")
	if clusterId == "" {
		c.JSON(400, gin.H{"error": "not exit clusterId"})
	}

	j, err := job.ClusterDelete(clusterId, user.(*User).Username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "deleteCluster success", "data": j})
}
