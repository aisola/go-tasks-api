package sqlite

import (
	"testing"
	"time"

	"example.com/tasks"
	"github.com/jmoiron/sqlx"
	"github.com/matryer/is"
)

func newInMemoryRepository(t *testing.T) *Repository {
	t.Helper()

	repo, err := New(":memory:")
	if err != nil {
		t.Fatalf("could not open in-memory sqlite3: %s", err)
	}

	return repo
}

func TestCreateTask(t *testing.T) {
	is := is.New(t)
	repo := newInMemoryRepository(t)

	err := repo.CreateTask(&tasks.Task{
		Text: "testing",
	})
	is.NoErr(err) // Error from CreateTask

	// TODO: write test for content
}

func TestRetrieveTask(t *testing.T) {
	is := is.New(t)
	repo := newInMemoryRepository(t)
	id := tasks.NewTaskID()

	_, err := repo.RetrieveTask(id)
	is.Equal(err, tasks.ErrTaskNotFound)

	sqlx.MustExec(repo.db,
		`INSERT INTO tasks (id, created_at, updated_at, text, is_complete) VALUES (?, ?, ?, ?, ?);`,
		id, time.Now().UTC(), time.Now().UTC(), "testing", false,
	)

	task, err := repo.RetrieveTask(id)
	is.NoErr(err)                  // Error from RetrieveTask
	is.Equal(id, task.ID)          // should be id
	is.Equal("testing", task.Text) // should be "testing"
}

func TestListTasks(t *testing.T) {
	is := is.New(t)
	repo := newInMemoryRepository(t)
	id := tasks.NewTaskID()
	sqlx.MustExec(repo.db,
		`INSERT INTO tasks (id, created_at, updated_at, text, is_complete) VALUES (?, ?, ?, ?, ?);`,
		id, time.Now().UTC(), time.Now().UTC(), "testing", false,
	)

	tasks, err := repo.ListTasks()
	is.NoErr(err)                      // Error from ListTask
	is.Equal(id, tasks[0].ID)          // should be id
	is.Equal("testing", tasks[0].Text) // should be "testing"
}

func TestUpdateTask(t *testing.T) {
	is := is.New(t)
	repo := newInMemoryRepository(t)
	id := tasks.NewTaskID()
	sqlx.MustExec(repo.db,
		`INSERT INTO tasks (id, created_at, updated_at, text, is_complete) VALUES (?, ?, ?, ?, ?);`,
		id, time.Now().UTC(), time.Now().UTC(), "changeme", false,
	)

	task, err := repo.UpdateTask(id, &tasks.Task{Text: "testing"})
	is.NoErr(err)                  // Error from UpdateTask
	is.Equal(id, task.ID)          // should be id
	is.Equal("testing", task.Text) // should be "testing"
}

func TestDeleteTask(t *testing.T) {
	is := is.New(t)
	repo := newInMemoryRepository(t)
	id := tasks.NewTaskID()

	is.NoErr(repo.DeleteTask(id)) // Error from DeleteTask

	sqlx.MustExec(repo.db,
		`INSERT INTO tasks (id, created_at, updated_at, text, is_complete) VALUES (?, ?, ?, ?, ?);`,
		id, time.Now().UTC(), time.Now().UTC(), "changeme", false,
	)

	is.NoErr(repo.DeleteTask(id)) // Error from DeleteTask
}
