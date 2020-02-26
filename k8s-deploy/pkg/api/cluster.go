package api

import (
	"fate-cloud-agent/pkg/db"
	"fate-cloud-agent/pkg/job"
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
		c.JSON(400, gin.H{"error": "err"})
	}
	cluster, err := db.ClusterFindByUUID(clusterId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{"data": cluster})
}

func (_ *Cluster) getClusterList(c *gin.Context) {

	clusterList, err := db.FindClusterList("")
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
		c.JSON(400, gin.H{"error": "err"})
	}

	j,err := job.ClusterDelete(clusterId, user.(*User).Username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "deleteCluster success", "data": j})
}
