package apperror

type Error struct {
	Code    int
	Message string
	Detail  any
}

func (e *Error) Error() string { return e.Message }

func New(code int, message string, detail any) *Error {
	return &Error{Code: code, Message: message, Detail: detail}
}
