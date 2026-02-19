package worker_constants

const (
	SUCCESS                  string = "SUCCESS"
	FAILED                   string = "FAILED"
	DispatchJobEndpoint      string = "%s/worker/jobs/dispatch"
	ReportCompletionEndpoint string = "%s/worker/jobs/%s/completion"
	HTTPWebhookExecutor      string = "http_webhook"
)
