package tasks

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	// ErrTaskNotFound is returned by repositories when a task is not found in
	// the respository.
	ErrTaskNotFound = errors.New("task not found")
)

// Task is the domain task implementation.
type Task struct {
	ID         string    `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Text       string    `db:"text"`
	IsComplete bool      `db:"is_complete"`
}

// TaskRepository defines the interface which repositories must implement in
// order to be used by the application.
type TaskRepository interface {
	CreateTask(t *Task) error
	ListTasks() ([]*Task, error)
	RetrieveTask(id string) (*Task, error)
	UpdateTask(id string, t *Task) (*Task, error)
	DeleteTask(id string) error
}

// NewTaskID creates a new task ID.
func NewTaskID() string {
	return uuid.Must(uuid.NewV4()).String()
}
