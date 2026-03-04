package database_constants

import (
	"fmt"
	"strings"
)

type JobExecutionStatus string

const (
	Running   JobExecutionStatus = "RUNNING"
	Todo      JobExecutionStatus = "TODO"
	Completed JobExecutionStatus = "COMPLETED"
	Error     JobExecutionStatus = "ERROR"
	Retry     JobExecutionStatus = "RETRY"
)

func ParseJobExecutionStatus(s string) (JobExecutionStatus, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case string(Running):
		return Running, nil
	case string(Todo):
		return Todo, nil
	case string(Completed):
		return Completed, nil
	case string(Retry):
		return Retry, nil
	case string(Error):
		return Error, nil
	default:
		return "", fmt.Errorf("invalid image type: %s", s)
	}
}
