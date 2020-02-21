package taskhttp

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func (h *Handler) tasksDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		id := chi.URLParam(r, "id")

		if err := h.repo.DeleteTask(id); err != nil {
			h.logger.Error("failed to delete task",
				zap.String("request_id", requestID),
				zap.String("task_id", id),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
