package taskhttp

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"example.com/tasks"
)

// Handler is an HTTP handler for the tasks API.
type Handler struct {
	router chi.Router
	logger *zap.Logger
	repo   tasks.TaskRepository
}

// New creates a new Handler
func New(logger *zap.Logger, tr tasks.TaskRepository) *Handler {
	h := &Handler{
		router: chi.NewRouter(),
		logger: logger,
		repo:   tr,
	}

	h.routes()

	return h
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
