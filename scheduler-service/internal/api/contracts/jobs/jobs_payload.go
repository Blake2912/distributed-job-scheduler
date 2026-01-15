package jobs

import "time"

type CreateJobsPayload struct {
	Name                    string    `json:"name" binding:"required"`
	Type                    string    `json:"type" binding:"required"`
	Config                  string    `json:"config" binding:"required"`
	TriggerTime             time.Time `json:"trigger_time" binding:"required"`
	ShouldRetryAfterBackoff *bool     `json:"should_retry_after_backoff" binding:"required"`
	IsRecurringJob          *bool     `json:"is_recurring_job" binding:"required"`
}
