package worker_constants

import "time"

const (
	SUCCESS                   string        = "SUCCESS"
	FAILED                    string        = "FAILED"
	LeaseJobExecutionEndpoint string        = "%s/api/worker/executions/lease"
	ReportCompletionEndpoint  string        = "%s/api/worker/executions/%s/complete"
	LeaseExtensionEndpoint    string        = "%s/api/worker/executions/%s/heartbeat"
	HTTPWebhookExecutor       string        = "http_webhook"
	LeaseDuration             time.Duration = 180 * time.Second
)
