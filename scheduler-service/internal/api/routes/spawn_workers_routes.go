package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerSpawnWorkerRoutes(r *gin.RouterGroup, c *container.Container) {
	worker := r.Group("/worker")

	worker.GET("spawnWorker", c.SpawnWorkersHandler.SpawnWorkers)
}
