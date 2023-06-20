package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mehix/go-todos/pkg/todos"
)

var addTodo = func(ctx context.Context, td todos.Todo) error {
	var payload bytes.Buffer
	if err := json.NewEncoder(&payload).Encode(td); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, *apiURL+"/todos/", &payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("todo not created. Response: %d", resp.StatusCode)
	}

	return nil
}

var totalTodos = func() (int, error) {
	req, err := http.NewRequest(http.MethodGet, *apiURL+"/todos/", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var all []todos.Todo
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return 0, err
	}

	return len(all), nil
}

var byTags = func(tags []string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, *apiURL+"/todos/search/tags", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-type", "application/json")
	q := req.URL.Query()
	q.Set("q", strings.Join(tags, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var all []todos.Todo
	if err := json.NewDecoder(resp.Body).Decode(&all); err != nil {
		return 0, err
	}

	return len(all), nil
}
