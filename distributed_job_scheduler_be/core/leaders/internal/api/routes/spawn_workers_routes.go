package routes

import (
	"example.com/leader/internal/container"
	"github.com/gin-gonic/gin"
)

func registerSpawnWorkerRoutes(r *gin.RouterGroup, c *container.Container) {
	worker := r.Group("/worker")

	worker.GET("spawnWorker", c.SpawnWorkersHandler.SpawnWorkers)
}
