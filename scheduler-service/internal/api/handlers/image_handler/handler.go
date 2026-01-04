package imagehandler

import (
	"net/http"

	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc imageservice.ImageService
}

func New(svc imageservice.ImageService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetByTypeAndVersion(c *gin.Context) {
	imageType := c.Query("type")
	version := c.Query("version")

	if imageType == "" || version == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "type and version are required",
		})
		return
	}

	parsedImageType, err := databaseconstants.ParseImageType(imageType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	img, err := h.svc.GetByTypeAndVersion(
		c.Request.Context(),
		parsedImageType,
		version,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "image not found",
		})
		return
	}

	c.JSON(http.StatusOK, toResponse(img))
}

func toResponse(img models.ImageInformation) gin.H {
	return gin.H{
		"Image":   img.Image,
		"Type":    img.ImageType,
		"Version": img.ImageVersion,
	}
}
