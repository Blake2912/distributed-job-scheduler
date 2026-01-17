package workerhealthcheckhandler

import (
	"net/http"

	workerhealthchecks "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/worker_health_checks"
	"github.com/gin-gonic/gin"
)

type WorkerHealthCheckHandler struct {
	svc workerhealthchecks.WorkerHealthChecks
}

func NewWorkerHealthCheckHander(svc workerhealthchecks.WorkerHealthChecks) *WorkerHealthCheckHandler {
	return &WorkerHealthCheckHandler{
		svc: svc,
	}
}

// SpawnWorkers godoc
// @Summary      Performs worker health checks
// @Description  For a given worker
// @Tags         HealthCheck
// @Accept       json
// @Produce      json
// @Param        workerId query string true "Worker Id"
// @Param        jobExecutionId query int true "Job Execution id" minimum(1)
// @Success      200 {object} map[string]string "ok"
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /healthCheck [get]
func (w *WorkerHealthCheckHandler) HealthCheck(c *gin.Context) {
	workerId := c.Query("workerId")
	jobExecutionId := c.Query("jobExecutionId")

	if workerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workerId is required field",
		})
		return
	}

	if jobExecutionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "jobExecutionId is required field",
		})
		return
	}

	err := w.svc.CheckHealth(c.Request.Context(), workerId, jobExecutionId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// SpawnWorkers godoc
// @Summary      Performs worker health checks
// @Description  For a given worker
// @Tags         HealthCheck
// @Accept       json
// @Produce      json
// @Param        workerId query string true "Worker Id"
// @Param        jobExecutionId query int true "Job Execution id" minimum(1)
// @Success      200 {object} map[string]string "ok"
// @Failure      400 {object} map[string]string "Invalid input"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /healthCheck [delete]
func (w *WorkerHealthCheckHandler) DeleteHealthCheck(c *gin.Context) {
	workerId := c.Query("workerId")
	jobExecutionId := c.Query("jobExecutionId")

	if workerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "workerId is required field",
		})
		return
	}

	if jobExecutionId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "jobExecutionId is required field",
		})
		return
	}

	err := w.svc.DeleteWorkerKeys(c.Request.Context(), workerId, jobExecutionId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
