package contracts

// JobToExecute represents job lease from scheduler
type JobToExecute struct {
	JobID   uint
	JobType string
	Payload string
}

type ReportCompletionRequest struct {
	Status    string
	Error     string
	Retryable bool
}
