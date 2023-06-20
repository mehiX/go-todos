package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/mehix/go-todos/pkg/todos"
)

var apiURL = flag.String("api-url", "http://127.0.0.1:7070", "Run tests against this URL")

func init() {
	flag.Parse()
}

func TestAddAndFetch(t *testing.T) {

	title := "some title"

	// prepare new Todo
	var payload bytes.Buffer
	if err := json.NewEncoder(&payload).Encode(todos.Todo{Title: title}); err != nil {
		t.Fatal(err)
	}

	// create Todo
	id := func() string {
		// prepare request to create new Todo
		req, err := http.NewRequest(http.MethodPost, *apiURL+"/todos/", &payload)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-type", "application/json")

		// send request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("wrong status code after create. expected: %d, got: %d", http.StatusCreated, resp.StatusCode)
		}

		// decode response from the server
		var newTodo todos.Todo
		if err := json.NewDecoder(resp.Body).Decode(&newTodo); err != nil {
			t.Fatal(err)
		}

		// validate response
		if newTodo.Title != title {
			t.Fatalf("wrong title. expected: %s, got: %s", title, newTodo.Title)
		}
		if _, err := uuid.Parse(newTodo.ID); err != nil {
			t.Fatalf("new ID is not a UUID. Error: %v", err)
		}

		return newTodo.ID
	}()

	// find by ID
	func() {
		// prepare find Todo by ID
		req, err := http.NewRequest(http.MethodGet, *apiURL+"/todos/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// decode response from the server
		var newTodo todos.Todo
		if err := json.NewDecoder(resp.Body).Decode(&newTodo); err != nil {
			t.Fatal(err)
		}

		// validate response
		if newTodo.Title != title {
			t.Fatalf("wrong title. expected: %s, got: %s", title, newTodo.Title)
		}
		if _, err := uuid.Parse(newTodo.ID); err != nil {
			t.Fatalf("new ID is not a UUID. Error: %v", err)
		}
	}()

	// delete
	func() {
		req, err := http.NewRequest(http.MethodDelete, *apiURL+"/todos/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		_, _ = io.Copy(io.Discard, resp.Body)

		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("wrong response to Delete. expected: %d, got: %d", http.StatusNoContent, resp.StatusCode)
		}
	}()

	// find missing by ID
	func() {
		// prepare find Todo by ID
		req, err := http.NewRequest(http.MethodGet, *apiURL+"/todos/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		_, _ = io.Copy(io.Discard, resp.Body)

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("wrong response for missing ID. expected: %d, got: %d", http.StatusNotFound, resp.StatusCode)
		}
	}()

}
