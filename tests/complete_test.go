package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mehix/go-todos/pkg/todos"
)

func TestCompleteTodo(t *testing.T) {

	ctx := context.Background()

	td, err := addTodo(ctx, todos.Todo{Title: "some todo", Tags: []string{"tag1", "tag2"}})
	if err != nil {
		t.Fatal(err)
	}

	if td.CompletedAt != nil {
		t.Fatalf("wrong completed_at for a new todo. got: %v\n", *td.CompletedAt)
	}

	req, err := http.NewRequest(http.MethodPost, *apiURL+"/todos/"+td.ID+"/complete", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var completed todos.Todo
	if err := json.NewDecoder(resp.Body).Decode(&completed); err != nil {
		t.Fatal(err)
	}

	if completed.CompletedAt == nil {
		t.Fatalf("todo not completed. CompletedAt is still nil")
	}
}
