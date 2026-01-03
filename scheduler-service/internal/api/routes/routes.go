package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	api := r.Group("/api")

	registerImageRoutes(api, c)
}
