package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
)

type HTTPSchedulerClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPSchedulerClient(
	baseURL string,
	httpClient *http.Client,
) SchedulerClient {
	return &HTTPSchedulerClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *HTTPSchedulerClient) LeaseJob(ctx context.Context) (*contracts.JobToExecute, error) {
	url := fmt.Sprintf(worker_constants.LeaseJobExecutionEndpoint, c.baseURL)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"lease job failed with status %d",
			resp.StatusCode,
		)
	}

	//decode response
	var parsed contracts.JobToExecute

	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return nil, err
	}

	fmt.Println(parsed.JobExecutionID)
	fmt.Println(parsed.JobType)
	fmt.Println(parsed.Payload)

	return &contracts.JobToExecute{
		JobExecutionID: parsed.JobExecutionID,
		JobType:        parsed.JobType,
		Payload:        parsed.Payload,
	}, nil
}

func (c *HTTPSchedulerClient) ReportCompletion(
	ctx context.Context,
	jobExecId uint,
	status string,
	executionError string,
	retryable bool,
) error {
	url := fmt.Sprintf(worker_constants.ReportCompletionEndpoint, c.baseURL, strconv.FormatUint(uint64(jobExecId), 10))

	body, err := json.Marshal(
		contracts.ReportCompletionRequest{
			Status:    status,
			Error:     executionError,
			Retryable: retryable,
		})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"report completion failed with status %d",
			resp.StatusCode,
		)
	}

	return nil
}
