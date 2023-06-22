package todo

// Represents an application level error.
type AppError struct {
	Code string `json:"code"`
}

func (e *AppError) Error() string {
	return "app error"
}

func (e *AppError) Status() int {
	return 400
}

func NewAppError(code string) error {
	return &AppError{Code: code}
}
