package api

import (
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
)

type Version struct {
}

// Router is cluster router definition method
func (c *Version) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	cluster := r.Group("/version")
	cluster.Use(authMiddleware.MiddlewareFunc())
	{
		cluster.GET("/", c.getVersion)
	}
}

func (_ *Version) getVersion(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "getVersion success", "version": service.GetVersion()})
}
