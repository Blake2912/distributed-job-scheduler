package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerJobsRoutes(r *gin.RouterGroup, c *container.Container) {
	jobs := r.Group("/jobs")

	jobs.POST("create", c.JobsHandler.CreateJobs)

}
