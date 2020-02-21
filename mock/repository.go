package mock

import (
	"sync"
	"time"

	"example.com/tasks"
)

// Repository is an in-memory implementation of a repository. This is safe for
// concurrent use so you may use it inside of an HTTP handler concurrently.
type Repository struct {
	mu   sync.RWMutex
	data map[string]*tasks.Task
}

// New creates a new Repository. Any tasks passed to the repository will be used
// to initialize the in-memory db.
func New(ts ...*tasks.Task) *Repository {
	data := make(map[string]*tasks.Task)

	for _, t := range ts {
		data[t.ID] = t
	}

	return &Repository{
		data: data,
	}
}

// CreateTask creates a new task. All fields except Task.Text will be
// overridden by defaults.
func (r *Repository) CreateTask(t *tasks.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t.ID = tasks.NewTaskID()
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	t.IsComplete = false

	r.data[t.ID] = t
	return nil
}

// ListTasks lists all tasks in the in-memory repo.
func (r *Repository) ListTasks() ([]*tasks.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*tasks.Task, 0)
	for _, t := range r.data {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// RetrieveTask retrieves the task from the repo by ID.
func (r *Repository) RetrieveTask(id string) (*tasks.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.data[id]
	if !ok {
		return nil, tasks.ErrTaskNotFound
	}

	return t, nil
}

// UpdateTask updates a task, by id, in the repo. If the task, does not exist,
// it will return tasks.ErrTaskNotFound. Only t.Text and t.IsCompleted are used
// to update the fields. The returned Task is the updated version of the task.
func (r *Repository) UpdateTask(id string, t *tasks.Task) (*tasks.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.data[id]
	if !ok {
		return nil, tasks.ErrTaskNotFound
	}

	var hasChanged bool

	// if t.Text has changed AND it is not empty
	if t.Text != e.Text && t.Text != "" {
		e.Text = t.Text
		hasChanged = true
	}

	if t.IsComplete != e.IsComplete {
		e.IsComplete = t.IsComplete
		hasChanged = true
	}

	if hasChanged {
		e.UpdatedAt = time.Now().UTC()
	}

	return e, nil
}

// DeleteTask deletes the task by ID. Attempting to delete a task with an ID
// which does not exist is not considered an error.
func (r *Repository) DeleteTask(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)

	return nil
}
