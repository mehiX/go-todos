package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mehix/go-todos/pkg/todos"
	"golang.org/x/sync/errgroup"
)

// workaround for: `flag provided but not defined: -test.paniconexit0` when passing arguments to the test binary
var _ = func() bool { testing.Init(); return true }()

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
		return fmt.Errorf("Todo not created. Response: %d", resp.StatusCode)
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

func TestAddMultipleConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skipf("skip long running test: %s\n", "TestAddMultipleConcurrent")
	}

	n := rand.Intn(100) + 50
	workers := 7
	maxRunningTime := 30 * time.Second

	fmt.Printf("Trying to create %d todos in a maximum of %v\n", workers*n, maxRunningTime)

	// get the initial count so that we don't have to empty the database before the test
	initialCount, err := totalTodos()
	if err != nil {
		t.Fatal(err)
	}

	maxRun, cancel := context.WithTimeout(context.Background(), maxRunningTime)
	defer cancel()

	g, gCtx := errgroup.WithContext(maxRun)

	for idx := 0; idx < workers; idx++ {
		fmt.Printf("Start worker %d for %d requests\n", idx+1, n)
		g.Go(func() error {
			for i := 0; i < n; i++ {
				select {
				case <-gCtx.Done():
					return gCtx.Err()
				case <-time.Tick(time.Duration(rand.Int63n(200)+100) * time.Millisecond):
					if err := addTodo(gCtx, todos.Todo{
						ID:    uuid.NewString(),
						Title: fmt.Sprintf("Todo number %d", i+1)}); err != nil {
						return err
					}
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil && err != context.DeadlineExceeded {
		t.Fatal(err)
	}

	finalCount, err := totalTodos()
	if err != nil {
		t.Fatal(err)
	}
	if initialCount+n*workers != finalCount {
		t.Fatalf("wrong number of todo's fetched. expected: %d, got: %d", initialCount+n*workers, finalCount)
	}
}
