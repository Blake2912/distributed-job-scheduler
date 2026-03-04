package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerWorkerRoutes(r *gin.RouterGroup, c *container.Container) {
	worker := r.Group("/worker")
	{
		executions := worker.Group("/executions")
		{
			executions.POST("/lease", c.WorkerHandler.LeaseNextJob)
			executions.POST("/:id/complete", c.WorkerHandler.ReportCompletion)
			executions.POST("/:id/heartbeat", c.WorkerHandler.ExtendLease)
		}
	}
}
