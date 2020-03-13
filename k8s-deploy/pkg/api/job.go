package api

import (
	"fate-cloud-agent/pkg/db"
	"github.com/gin-gonic/gin"
)

type Job struct {
}

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

func (_ *Job) getJobList(c *gin.Context) {

	jobList, err := db.JobFindList("")
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"data": jobList, "msg": "getJobList success"})
}

func (_ *Job) getJob(c *gin.Context) {

	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, gin.H{"error": "not exit jobId"})
		return
	}
	result, err := db.JobFindByUUID(jobId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"data": result})
}

func (_ *Job) deleteJob(c *gin.Context) {

	jobId := c.Param("jobId")
	if jobId == "" {
		c.JSON(400, gin.H{"error": "not exit jobId"})
		return
	}

	err := db.JobDeleteByUUID(jobId)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"msg": "delete Job success"})
}
