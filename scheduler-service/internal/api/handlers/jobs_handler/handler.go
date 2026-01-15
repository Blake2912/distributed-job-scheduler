package jobshandler

import (
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/contracts/jobs"
	jobsservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/jobs_service"
	"github.com/gin-gonic/gin"
)

type JobsHandler struct {
	svc jobsservice.JobsService
}

func New(svc jobsservice.JobsService) *JobsHandler {
	return &JobsHandler{svc: svc}
}

// CreateJobs godoc
// @Summary      Create multiple jobs
// @Description  Creates multiple jobs in a single request (bulk create)
// @Tags         jobs
// @Accept       json
// @Produce      json
// @Param        jobs body []jobs.CreateJobsPayload true "List of jobs to create"
// @Success      201 {object} map[string]string "Jobs created successfully"
// @Failure      400 {object} map[string]string "Invalid request payload"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /jobs/create [post]`
func (h *JobsHandler) CreateJobs(c *gin.Context) {
	var req []jobs.CreateJobsPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.svc.CreateJobs(c, req)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "jobs were successfully created",
	})
}
