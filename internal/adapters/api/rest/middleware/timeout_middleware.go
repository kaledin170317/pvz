package middleware

import (
	"context"
	"net/http"
	"time"
)

const StartTimeKey string = "startTime"

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			ctx = context.WithValue(ctx, StartTimeKey, time.Now())
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
