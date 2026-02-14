package worker_constants

const (
	SUCCESS                  string = "SUCCESS"
	FAILED                   string = "FAILED"
	LeaseEndpoint            string = "%s/jobs/lease"
	ReportCompletionEndpoint string = "%s/jobs/%s/completion"
	HTTPWebhookExecutor      string = "http_webhook"
)
