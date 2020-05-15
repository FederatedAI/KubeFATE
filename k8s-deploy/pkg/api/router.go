package api

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
)

const ApiVersion = "v1"

func Router(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg":             "kubefate run success",
			"api_version":     ApiVersion,
			"service_version": service.GetVersion(),
		})
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "error path"})
		})

		cluster := new(Cluster)
		cluster.Router(v1)

		user := new(User)
		user.Router(v1)

		job := new(Job)
		job.Router(v1)

		chart := new(Chart)
		chart.Router(v1)

		version := new(Version)
		version.Router(v1)
	}
}
