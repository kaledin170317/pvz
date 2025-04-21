package middleware

import (
	"net/http"
	"pvZ/internal/metrics"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		path := r.URL.Path
		method := r.Method

		metrics.HTTPRequestTotal.WithLabelValues(path, method).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(path, method).Observe(duration.Seconds())
	})
}
