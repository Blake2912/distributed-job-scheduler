package worker_constants

import "time"

const (
	SUCCESS                  string        = "SUCCESS"
	FAILED                   string        = "FAILED"
	DispatchJobEndpoint      string        = "%s/worker/jobs/lease"
	ReportCompletionEndpoint string        = "%s/worker/jobs/%s/completion"
	HTTPWebhookExecutor      string        = "http_webhook"
	LeaseDuration            time.Duration = 180 * time.Second
)
