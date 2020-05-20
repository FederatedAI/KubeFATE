package api

import (
	"fmt"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/db"
	"github.com/FederatedAI/KubeFATE/k8s-deploy/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type Chart struct {
}

// Router is cluster router definition method
func (c *Chart) Router(r *gin.RouterGroup) {

	authMiddleware, _ := GetAuthMiddleware()
	cluster := r.Group("/chart")
	cluster.Use(authMiddleware.MiddlewareFunc())
	{
		cluster.POST("", c.createChart)
		cluster.GET("/", c.getChartList)
		cluster.GET("/:chartId", c.getChart)
		cluster.DELETE("/:chartId", c.deleteChart)
	}
}

func (_ *Chart) createChart(c *gin.Context) {
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chartRequested, err := loader.LoadArchive(f)

	helmChart, err := service.ChartRequestedTohelmChart(chartRequested)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	helmUUID, err := db.ChartSave(helmChart)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "createChart success", "data": gin.H{"helmUUID": helmUUID}})
}

func (_ *Chart) getChartList(c *gin.Context) {

	chartList, err := db.FindHelmChartList()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "getChartList success", "data": chartList})
}

func (_ *Chart) getChart(c *gin.Context) {

	chartId := c.Param("chartId")
	if chartId == "" {
		c.JSON(400, gin.H{"error": "not exit chartId"})
		return
	}

	chartList, err := db.FindHelmChart(chartId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "getChart success", "data": chartList})
}

func (_ *Chart) deleteChart(c *gin.Context) {
	chartId := c.Param("chartId")
	if chartId == "" {
		c.JSON(400, gin.H{"error": "not exit chartId"})
		return
	}
	chart := new(db.HelmChart)
	n, err := db.DeleteByUUID(chart, chartId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if n != 1 {
		c.JSON(200, gin.H{"msg": "deleteChart error, DeletedCount=" + fmt.Sprintf("%d", n)})
		return
	}
	c.JSON(200, gin.H{"msg": "deleteChart success"})
}
