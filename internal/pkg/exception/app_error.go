package exception

import "google.golang.org/grpc/codes"

// AppError handles application exception.
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Status  codes.Code             `json:"status"`
	Details map[string]interface{} `json:"details"`
}

func (e AppError) Error() string {
	return e.Message
}

// New functions create a new AppError instance
func New(code, message string) AppError {
	return AppError{Code: code, Message: message}
}
