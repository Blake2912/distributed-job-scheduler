package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerWorkerHealthCheckRoutes(r *gin.RouterGroup, c *container.Container) {
	worker := r.Group("/healthCheck")

	worker.GET("", c.WorkerHealthCheckHandler.HealthCheck)
	worker.DELETE("", c.WorkerHealthCheckHandler.DeleteHealthCheck)
}
