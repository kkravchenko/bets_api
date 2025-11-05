package errors

type ValidationError struct {
	Msg  string
	Code int
}

func (v *ValidationError) Error() string {
	return v.Msg
}

func NewValidationError(msg string, code int) error {
	return &ValidationError{Msg: msg, Code: code}
}
