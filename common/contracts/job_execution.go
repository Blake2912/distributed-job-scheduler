package contracts

// JobToExecute represents job lease from scheduler
type JobToExecute struct {
	JobExecutionID uint
	JobType        string
	Payload        string
}

type ReportCompletionRequest struct {
	Status    string
	Error     string
	Retryable bool
}
