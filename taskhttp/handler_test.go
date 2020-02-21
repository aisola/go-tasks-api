package taskhttp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"

	"example.com/tasks"
	"example.com/tasks/mock"
)

func callWithNewHandler(t *testing.T, req *http.Request, ts ...*tasks.Task) *httptest.ResponseRecorder {
	// Mark this as a testing helper function.
	t.Helper()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to
	// record the response.
	rr := httptest.NewRecorder()

	// Instantiate a new handler and call the ServeHTTP method to simulate an
	// HTTP request.
	New(zap.NewNop(), mock.New(ts...)).ServeHTTP(rr, req)

	// Return the ResponseRecorder so that our real tests can do their thing.
	return rr
}

func TestTasksCreate(t *testing.T) {
	is := is.New(t)

	data := bytes.NewBuffer([]byte(`{"text": "testing"}`))

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodPost, "/", data)
	if err != nil {
		t.Fatal(err)
	}

	rr := callWithNewHandler(t, req)
	is.Equal(rr.Code, http.StatusCreated)                           // Status should equal 201
	is.True(strings.Contains(rr.Body.String(), `"text":"testing"`)) // Body -> text = testing
}

func TestTasksUpdate(t *testing.T) {
	is := is.New(t)

	data := bytes.NewBuffer([]byte(`{"text": "testing"}`))
	id := tasks.NewTaskID()

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodPatch, "/"+id, data)
	if err != nil {
		t.Fatal(err)
	}

	rr := callWithNewHandler(t, req, &tasks.Task{
		ID:         id,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Text:       "changeme",
		IsComplete: false,
	})
	is.Equal(rr.Code, http.StatusOK)                                  // Status should equal 200
	is.True(strings.Contains(rr.Body.String(), `"id":"`+id+`"`))      // Body -> id is our id
	is.True(!strings.Contains(rr.Body.String(), `"text":"changeme"`)) // Body -> text != changeme
	is.True(strings.Contains(rr.Body.String(), `"text":"testing"`))   // Body -> text = testing
}

func TestTasksRetrieve(t *testing.T) {
	is := is.New(t)

	id := tasks.NewTaskID()

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, "/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := callWithNewHandler(t, req)
	is.Equal(rr.Code, http.StatusNotFound) // Status should equal 404

	rr = callWithNewHandler(t, req, &tasks.Task{
		ID:         id,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Text:       "testing",
		IsComplete: false,
	})
	is.Equal(rr.Code, http.StatusOK)                             // Status should equal 200
	is.True(strings.Contains(rr.Body.String(), `"id":"`+id+`"`)) // Body -> id is our id
}

func TestTasksList(t *testing.T) {
	is := is.New(t)

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	id1 := tasks.NewTaskID()
	id2 := tasks.NewTaskID()

	rr := callWithNewHandler(t, req, &tasks.Task{
		ID:         id1,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Text:       "testing",
		IsComplete: false,
	}, &tasks.Task{
		ID:         id2,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Text:       "testing",
		IsComplete: false,
	})
	is.Equal(rr.Code, http.StatusOK)                              // Status should equal 200
	is.True(strings.Contains(rr.Body.String(), `"id":"`+id1+`"`)) // Body -> id is our id1
	is.True(strings.Contains(rr.Body.String(), `"id":"`+id2+`"`)) // Body -> id is our id2
}

func TestTasksDelete(t *testing.T) {
	is := is.New(t)

	id := tasks.NewTaskID()

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodDelete, "/"+id, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := callWithNewHandler(t, req)
	is.Equal(rr.Code, http.StatusNoContent) // Status should equal 204
	is.Equal(len(rr.Body.String()), 0)      // Non-empty response body

	rr = callWithNewHandler(t, req, &tasks.Task{
		ID:         id,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Text:       "testing",
		IsComplete: false,
	})
	is.Equal(rr.Code, http.StatusNoContent) // Status should equal 204
	is.Equal(len(rr.Body.String()), 0)      // Non-empty response body
}
