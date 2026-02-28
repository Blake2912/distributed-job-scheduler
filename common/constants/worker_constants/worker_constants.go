package worker_constants

import "time"

const (
	SUCCESS                   string        = "SUCCESS"
	FAILED                    string        = "FAILED"
	LeaseJobExecutionEndpoint string        = "%s/worker/executions/lease"
	ReportCompletionEndpoint  string        = "%s/worker/executions/%s/complete"
	LeaseExtensionEndpoint    string        = "%s/worker/executions/%s/heartbeat"
	HTTPWebhookExecutor       string        = "http_webhook"
	LeaseDuration             time.Duration = 180 * time.Second
)
