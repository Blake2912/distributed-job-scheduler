package workerhandler

import (
	"net/http"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/worker"
	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	svc worker.WorkerJobDispatchService
}

func NewWorkerHandler(svc worker.WorkerJobDispatchService) *WorkerHandler {
	return &WorkerHandler{
		svc: svc,
	}
}

// DispatchNext godoc
// @Summary      Dispatch next job
// @Description  Provides next job to worker for execution
// @Tags         worker
// @Accept       json
// @Produce      json
// @Param        jobs body []jobs.CreateJobsPayload true "List of jobs to create"
// @Success      201 {object} contracts.JobToExecute
// @Failure      400 {object} map[string]string "Invalid request payload"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /worker/jobs/dispatch [get]`
func (h *WorkerHandler) DispatchNext(c *gin.Context) {
	jobToExecute, err := h.svc.DispatchNextJob(c.Request.Context())
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
