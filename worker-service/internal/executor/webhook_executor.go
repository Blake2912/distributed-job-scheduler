package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WebhookExecutor struct {
	httpClient *http.Client
}

func NewWebhookExecutor(client *http.Client) *WebhookExecutor {
	return &WebhookExecutor{
		httpClient: client,
	}
}

func (executor *WebhookExecutor) Execute(
	ctx context.Context,
	payload string,
) error {
	var data HTTPWebhookPayload

	//deserialize payload
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	//validate payload
	if data.URL == "" {
		return fmt.Errorf("url is required")
	}

	if data.Method == "" {
		//default to POST
		data.Method = http.MethodPost
	}

	if data.TimeoutSeconds <= 0 {
		//default timeout
		data.TimeoutSeconds = 10
	}

	reqCtx, cancel := context.WithTimeout(
		ctx,
		time.Duration(data.TimeoutSeconds)*time.Second,
	)

	defer cancel()

	//building HTTP request
	req, err := http.NewRequestWithContext(
		reqCtx,
		strings.ToUpper(data.Method),
		data.URL,
		bytes.NewBuffer([]byte(data.Body)),
	)

	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	//set headers
	for k, v := range data.Headers {
		req.Header.Set(k, v)
	}

	//Idempotency header (this is required for deduplication)
	req.Header.Set("Idempotency-Key", uuid.NewString())

	//execution
	resp, err := executor.httpClient.Do(req)
	if err != nil {
		return NewRetryable(fmt.Errorf("http request failed: %w", err))
	}
	defer resp.Body.Close()

	//classify response for retry based on statuscode
	if resp.StatusCode >= 500 {
		return NewRetryable(
			fmt.Errorf("server error: %d", resp.StatusCode),
		)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("client error: %d", resp.StatusCode)
	}

	return nil
}
