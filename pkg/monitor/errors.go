package monitor

type RequestError struct {
	URL     string
	Message string
}

// Error returns the RequestError Message.
func (e *RequestError) Error() string {
	return e.Message
}
