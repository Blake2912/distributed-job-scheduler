package httpclient

import "fmt"

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.StatusCode, e.Body)
}

func NewHTTPError(code int, body []byte) error {
	return &HTTPError{
		StatusCode: code,
		Body:       string(body),
	}
}
