package error

type BaseError struct {
	Code    string
	Message string
}

func NewError(code string, message string) BaseError {
	return BaseError{
		Code:    code,
		Message: message,
	}
}

func (e BaseError) Error() string {
	return e.Message
}
