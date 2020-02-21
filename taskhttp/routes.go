package taskhttp

import (
	"time"

	"github.com/go-chi/chi/middleware"
)

func (h *Handler) routes() {
	// A good base middleware stack
	h.router.Use(middleware.RequestID)
	h.router.Use(middleware.RealIP)
	h.router.Use(logAccess(h.logger.Named("access")))
	h.router.Use(middleware.Recoverer)

	// TODO: probably should make this timeout configurable
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	h.router.Use(middleware.Timeout(2 * time.Second))

	h.router.Get("/", h.tasksList())
	h.router.Post("/", h.tasksCreate())
	h.router.Get("/{id}", h.tasksRetrieve())
	h.router.Patch("/{id}", h.tasksUpdate())
	h.router.Delete("/{id}", h.tasksDelete())
}
