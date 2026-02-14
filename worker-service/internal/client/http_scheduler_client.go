package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
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

func (c *HTTPSchedulerClient) LeaseJob(ctx context.Context) (*LeaseJobResponse, error) {
	url := fmt.Sprintf(worker_constants.LeaseEndpoint, c.baseURL)

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
	var parsed LeaseJobResponse

	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return nil, err
	}

	fmt.Println(parsed.JobID)
	fmt.Println(parsed.JobType)
	fmt.Println(parsed.Payload)

	return &LeaseJobResponse{
		JobID:   parsed.JobID,
		JobType: parsed.JobType,
		Payload: parsed.Payload,
	}, nil
}

func (c *HTTPSchedulerClient) ReportCompletion(
	ctx context.Context,
	jobId string,
	status string,
	executionError string,
	retryable bool,
) error {
	url := fmt.Sprintf(worker_constants.ReportCompletionEndpoint, c.baseURL, jobId)

	body, err := json.Marshal(
		ReportCompletionRequest{
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
