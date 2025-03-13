package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Use a distinct function name to avoid clashes
func LogrusLogger(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"ip":       r.RemoteAddr,
				"duration": time.Since(start).Milliseconds(),
				"status":   w.Header().Get("Status"), // Extract status if available
			}).Info("Request started")

			next.ServeHTTP(w, r)
		})
	}
}
