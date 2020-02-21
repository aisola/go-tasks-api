package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"example.com/tasks"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const initializeTableQuery = `
CREATE TABLE IF NOT EXISTS tasks (
	id TEXT,
	created_at DATETIME,
	updated_at DATETIME,
	text TEXT,
	is_complete BOOLEAN
);
`

// Repository is an sqlite3 implementation of a repository.
type Repository struct {
	db *sqlx.DB
}

// New connects to a database, creating it if it doesn't exist, and
// initializes a repository. Errors come from connection issues. This
// method will panic if it encounters an error setting up the tables in the
// database at all.
func New(s string) (*Repository, error) {
	db, err := sqlx.Connect("sqlite3", s)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	sqlx.MustExec(db, initializeTableQuery)

	repo := &Repository{
		db: db,
	}

	return repo, nil
}

// CreateTask creates a new task. All fields except Task.Text will be
// overridden by defaults.
func (r *Repository) CreateTask(t *tasks.Task) error {
	const query = "INSERT INTO tasks (id, created_at, updated_at, text, is_complete) VALUES (?, ?, ?, ?, ?);"

	t.ID = tasks.NewTaskID()
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	t.IsComplete = false

	if _, err := r.db.Exec(query, t.ID, t.CreatedAt, t.UpdatedAt, t.Text, t.IsComplete); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// ListTasks lists all tasks in the repo.
func (r *Repository) ListTasks() ([]*tasks.Task, error) {
	const query = "SELECT * FROM tasks;"
	ts := make([]*tasks.Task, 0)

	if err := r.db.Select(&ts, query); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return ts, nil
}

// RetrieveTask retrieves the task from the repo by ID.
func (r *Repository) RetrieveTask(id string) (*tasks.Task, error) {
	const query = "SELECT * FROM tasks WHERE id=? LIMIT 1;"
	task := &tasks.Task{}

	if err := r.db.Get(task, query, id); err == sql.ErrNoRows {
		return nil, tasks.ErrTaskNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve task: %w", err)
	}

	return task, nil
}

// UpdateTask updates a task, by id, in the repo. If the task, does not exist,
// it will return tasks.ErrTaskNotFound. Only t.Text and t.IsCompleted are used
// to update the fields. The returned Task is the updated version of the task.
func (r *Repository) UpdateTask(id string, t *tasks.Task) (*tasks.Task, error) {
	const query = "UPDATE tasks SET text=?, is_complete=? WHERE id=?;"
	const getQuery = "SELECT * FROM tasks WHERE id=? LIMIT 1;"
	var task tasks.Task

	if _, err := r.db.Exec(query, t.Text, t.IsComplete, id); err == sql.ErrNoRows {
		return nil, tasks.ErrTaskNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	if err := r.db.Get(&task, getQuery, id); err != nil {
		return nil, fmt.Errorf("failed to retrieve task after update: %w", err)
	}

	return &task, nil
}

// DeleteTask deletes the task by ID. Attempting to delete a task with an ID
// which does not exist is not considered an error.
func (r *Repository) DeleteTask(id string) error {
	const query = "DELETE FROM tasks WHERE id=?;"

	if _, err := r.db.Exec(query, id); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
