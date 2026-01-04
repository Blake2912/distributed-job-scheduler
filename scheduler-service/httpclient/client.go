package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	http *http.Client
}

func New(timeout time.Duration) *Client {
	return &Client{
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

// Generic Function to make HTTP Calls to other services
// Future Improvements: Add retries for these calls when it fails
func Do[T any](
	ctx context.Context,
	c *Client,
	reqCfg Request,
) (T, error) {

	var zero T

	var body io.Reader
	if reqCfg.Body != nil {
		b, err := json.Marshal(reqCfg.Body)
		if err != nil {
			return zero, fmt.Errorf("marshal request body: %w", err)
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		reqCfg.Method,
		reqCfg.URL,
		body,
	)
	if err != nil {
		return zero, err
	}

	// headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range reqCfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, NewHTTPError(resp.StatusCode, respBody)
	}

	var out T
	if len(respBody) == 0 {
		return out, nil
	}

	if err := json.Unmarshal(respBody, &out); err != nil {
		return zero, fmt.Errorf("unmarshal response: %w", err)
	}

	return out, nil
}
