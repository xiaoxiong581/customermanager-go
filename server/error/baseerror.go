package error

type BaseError struct {
	Code    string
	Message string
}

func (e BaseError) Error() string {
	return e.Message
}
