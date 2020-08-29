package exception

// AppError handles application exception.
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

func (e AppError) Error() string {
	return e.Message
}

// New functions create a new AppError instance
func New(code, message string) AppError {
	return AppError{Code: code, Message: message}
}
