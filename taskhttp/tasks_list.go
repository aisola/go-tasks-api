package taskhttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

func (h *Handler) tasksList() http.HandlerFunc {
	type responseTask struct {
		ID         string    `json:"id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		Text       string    `json:"text"`
		IsComplete bool      `json:"is_complete"`
	}
	type response struct {
		Length int             `json:"length"`
		Items  []*responseTask `json:"items"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		ts, err := h.repo.ListTasks()
		if err != nil {
			h.logger.Error("failed to find task",
				zap.String("request_id", requestID),
				zap.Error(err),
			)
			internalServerError(w)
			return
		}

		l := len(ts)
		res := &response{
			Length: l,
			Items:  make([]*responseTask, l),
		}

		for i, t := range ts {
			res.Items[i] = &responseTask{
				ID:         t.ID,
				CreatedAt:  t.CreatedAt,
				UpdatedAt:  t.UpdatedAt,
				Text:       t.Text,
				IsComplete: t.IsComplete,
			}
		}

		respondJSON(w, http.StatusOK, res)
	}
}
