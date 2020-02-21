package taskhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"

	"example.com/tasks"
)

func (h *Handler) tasksUpdate() http.HandlerFunc {
	type request struct {
		Text       string `json:"text"`
		IsComplete bool   `json:"is_complete"`
	}
	type response struct {
		ID         string    `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Text       string    `json:"text"`
		IsComplete bool      `json:"is_complete"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			requestID = middleware.GetReqID(r.Context())
			id        = chi.URLParam(r, "id")
			req       request
		)
		if err := decode(r, &req); err != nil {
			h.logger.Error("failed to decode request",
				zap.String("request_id", requestID),
				zap.String("task_id", id),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		task := &tasks.Task{
			Text:       req.Text,
			IsComplete: req.IsComplete,
		}

		task, err := h.repo.UpdateTask(id, task)
		if err == tasks.ErrTaskNotFound {
			h.logger.Warn("task not found",
				zap.String("request_id", requestID),
				zap.String("task_id", id),
				zap.Error(err),
			)
			respondJSONError(w, http.StatusNotFound, "task not found")
			return
		} else if err != nil {
			h.logger.Error("failed to update task",
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
