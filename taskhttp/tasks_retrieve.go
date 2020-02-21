package taskhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"example.com/tasks"
)

func (h *Handler) tasksRetrieve() http.HandlerFunc {
	type response struct {
		ID         string    `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Text       string    `json:"text"`
		IsComplete bool      `json:"is_complete"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		id := chi.URLParam(r, "id")

		task, err := h.repo.RetrieveTask(id)
		if err == tasks.ErrTaskNotFound {
			h.logger.Warn("task not found",
				zap.String("request_id", requestID),
				zap.String("task_id", id),
				zap.Error(err),
			)
			respondJSONError(w, http.StatusNotFound, "task not found")
			return
		} else if err != nil {
			h.logger.Error("failed to find task",
				zap.String("request_id", requestID),
				zap.String("task_id", id),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		respondJSON(w, http.StatusOK, &response{
			ID:         task.ID,
			CreatedAt:  task.CreatedAt,
			UpdatedAt:  task.UpdatedAt,
			Text:       task.Text,
			IsComplete: task.IsComplete,
		})
	}
}
