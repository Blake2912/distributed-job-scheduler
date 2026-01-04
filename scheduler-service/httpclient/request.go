package httpclient

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    any
}
