package errs

import (
	"context"
	"net/http"

	"gitlab.com/romalor/roxi"

	"gitlab.com/romalor/htmx-contacts/middleware/logging"
)

func Errors(log logging.Logger) roxi.MiddlewareFunc {
	return func(next roxi.HandlerFunc) roxi.HandlerFunc {
		return func(ctx context.Context, r *http.Request) error {
			if err := next(ctx, r); err != nil {
				log("caught error in request",
					"error", err)
				return err
			}
			return nil
		}
	}
}
