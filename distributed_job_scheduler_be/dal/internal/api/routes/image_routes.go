package routes

import (
	"example.com/dal/internal/container"
	"github.com/gin-gonic/gin"
)

func registerImageRoutes(r *gin.RouterGroup, c *container.Container) {
	images := r.Group("/images")

	images.GET("getByTypeAndVersion", c.ImageHandler.GetByTypeAndVersion)
}
