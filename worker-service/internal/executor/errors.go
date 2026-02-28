package executor

//creating this new type which implements Error() interface so that we can easily use it to identify retryable error
//by calling the IsRetryable() method
type RetryableError struct {
	Err error
}

func (r RetryableError) Error() string {
	return r.Err.Error()
}

//constructor
func NewRetryable(err error) error {
	return RetryableError{Err: err}
}

func IsRetryable(err error) bool {
	_, ok := err.(RetryableError)
	return ok
}
