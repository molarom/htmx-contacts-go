package logging

import (
	"context"
	"net/http"
	"time"

	"gitlab.com/romalor/roxi"
)

type Logger func(msg string, args ...any)

func Logging(log Logger) roxi.MiddlewareFunc {
	return func(next roxi.HandlerFunc) roxi.HandlerFunc {
		return func(ctx context.Context, r *http.Request) error {
			before := time.Now()

			w := &statusWriter{ResponseWriter: roxi.GetWriter(ctx)}
			next.ServeHTTP(w, r.WithContext(ctx))

			log(
				"handled request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", w.status,
				"duration", time.Since(before),
			)

			return nil
		}
	}
}

type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (rw *statusWriter) WriteHeader(code int) {
	if rw.written {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.written = true
}

func (rw *statusWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}
