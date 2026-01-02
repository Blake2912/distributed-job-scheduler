package routes

import (
	"example.com/dal/internal/container"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	api := r.Group("/api")

	registerImageRoutes(api, c)
}
