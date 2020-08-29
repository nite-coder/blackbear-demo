package gql

import "github.com/jasonsoft/starter/internal/pkg/exception"

var (
	ErrNotFound = exception.New("NOT_FOUND", "resource was not found")
)
