package taskhttp

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func logAccess(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestID := middleware.GetReqID(r.Context())
			metrics := httpsnoop.CaptureMetrics(next, w, r)

			logger.Info("request completed",
				zap.String("request_id", requestID),
				zap.String("request_method", r.Method),
				zap.String("request_path", r.URL.Path),
				zap.Duration("request_duration", metrics.Duration),
				zap.Int("response_code", metrics.Code),
				zap.Int64("response_length", metrics.Written),
			)

		})
	}
}
