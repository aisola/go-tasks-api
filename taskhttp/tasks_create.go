package taskhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"example.com/tasks"
)

func (h *Handler) tasksCreate() http.HandlerFunc {
	type request struct {
		Text string `json:"text"`
	}
	type response struct {
		ID         string    `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Text       string    `json:"text"`
		IsComplete bool      `json:"is_complete"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		var req request
		if err := decode(r, &req); err != nil {
			h.logger.Error("failed to decode request",
				zap.String("request_id", requestID),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		task := &tasks.Task{
			Text: req.Text,
		}

		if err := h.repo.CreateTask(task); err != nil {
			h.logger.Error("failed to create task",
				zap.String("request_id", requestID),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		respondJSON(w, http.StatusCreated, &response{
			ID:         task.ID,
			CreatedAt:  task.CreatedAt,
			UpdatedAt:  task.UpdatedAt,
			Text:       task.Text,
			IsComplete: task.IsComplete,
		})
	}
}
