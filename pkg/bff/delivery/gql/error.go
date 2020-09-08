package gql

import (
	"context"
	"errors"
	"fmt"

	"github.com/99designs/gqlgen/handler"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/pkg/domain"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.temporal.io/sdk/temporal"
)

var (
	ErrNotFound = domain.New("NOT_FOUND", "resource was not found")
)

func CentralizedError() handler.Option {
	return handler.ErrorPresenter(
		func(ctx context.Context, err error) *gqlerror.Error {
			logger := log.FromContext(ctx)

			var appErr domain.AppError
			if errors.As(err, &appErr) {
				gErr := &gqlerror.Error{
					Message: appErr.Error(),
					Extensions: map[string]interface{}{
						"code": appErr.Code,
					},
				}
				return gErr
			}

			// workflow
			var workflowErr *temporal.ApplicationError
			if errors.As(err, &workflowErr) {
				if workflowErr.Type() != "" {
					gErr := &gqlerror.Error{
						Message: workflowErr.Error(),
						Extensions: map[string]interface{}{
							"code": workflowErr.Type(),
						},
					}
					return gErr
				}

				myErr := workflowErr.Unwrap()
				if myErr != nil {
					err = myErr
				}

			}

			// unknow error
			gErr := &gqlerror.Error{
				Message: err.Error(),
				Extensions: map[string]interface{}{
					"code": "UNKNOWN_ERROR",
				},
			}

			logger.Err(err).Error("gql: unknown error")
			return gErr
		})
}

func RecoverError() handler.Option {
	return handler.RecoverFunc(func(ctx context.Context, err interface{}) error {
		logger := log.FromContext(ctx)
		myErr := fmt.Errorf("Internal server error! %w", err)

		logger.Err(myErr).Error("Internal server error!")
		return myErr
	})
}
