package workerhandler

import (
	"net/http"
	"strconv"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/worker"
	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	svc worker.WorkerLeaseJobService
}

func NewWorkerHandler(svc worker.WorkerLeaseJobService) *WorkerHandler {
	return &WorkerHandler{
		svc: svc,
	}
}

// LeaseNextJob godoc
// @Summary      Leases next job
// @Description  Provides next job to worker for execution
// @Tags         worker
// @Accept       json
// @Produce      json
// @Success      200 {object} contracts.JobToExecute
// @Failure      400 {object} map[string]string "Invalid request payload"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /worker/executions/lease [get]`
func (h *WorkerHandler) LeaseNextJob(c *gin.Context) {
	jobToExecute, err := h.svc.LeaseJob(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if jobToExecute == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, jobToExecute)
}

// ReportCompletion godoc
// @Summary      Report execution completion
// @Description  Provides the information related to job execution completion from worker
// @Tags         worker
// @Accept       json
// @Produce      json
// @Param        jobs body contracts.ReportCompletionRequest true "execution completion info"
// @Success      200
// @Failure      400 {object} map[string]string "Invalid request payload"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /worker/executions/:id/complete [post]`
func (h *WorkerHandler) ReportCompletion(c *gin.Context) {
	var req contracts.ReportCompletionRequest

	id := c.Param("id")
	execId, err := strconv.ParseUint(id, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid execution id",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.svc.CompleteJobExecution(c.Request.Context(), uint(execId), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
