package spawnworkers

import (
	"net/http"
	"strconv"

	spawnworkers "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/spawn_workers"
	"github.com/gin-gonic/gin"
)

type SpawnWorkersHandler struct {
	svc spawnworkers.SpawnWorkersService
}

func NewSpawnWokersHandler(svc spawnworkers.SpawnWorkersService) *SpawnWorkersHandler {
	return &SpawnWorkersHandler{svc: svc}
}

// SpawnWorkers godoc
// @Summary      Spawn worker pods
// @Description  Spawns the given number of worker pods in the cluster
// @Tags         worker
// @Accept       json
// @Produce      json
// @Param        noOfWorkers query int true "Number of workers to spawn" minimum(1)
// @Success      200 {object} map[string]string "Pods spawned successfully"
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /worker/spawnWorker [get]
func (h *SpawnWorkersHandler) SpawnWorkers(c *gin.Context) {
	noOfWorkers := c.Query("noOfWorkers")

	if noOfWorkers == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "noOfWorkers are required",
		})
		return
	}

	parsedNoOfWorkers, err := strconv.Atoi(noOfWorkers)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "noOfWorkers must be a valid integer",
		})
		return
	}

	if parsedNoOfWorkers <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "noOfOWorkers must be greater than 0",
		})
		return
	}

	if err := h.svc.SpawnWorkers(c.Request.Context(), parsedNoOfWorkers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pods spawnned succesfully",
	})

}
