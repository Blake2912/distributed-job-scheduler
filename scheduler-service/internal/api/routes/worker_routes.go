package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerWorkerRoutes(r *gin.RouterGroup, c *container.Container) {
	worker := r.Group("/worker")

	worker.GET("jobs/dispatch", c.WorkerHandler.DispatchNext)
}
