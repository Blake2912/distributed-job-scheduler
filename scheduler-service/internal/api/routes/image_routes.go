package routes

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/container"
	"github.com/gin-gonic/gin"
)

func registerImageRoutes(r *gin.RouterGroup, c *container.Container) {
	images := r.Group("/images")

	images.GET("getByTypeAndVersion", c.ImageHandler.GetByTypeAndVersion)
}
