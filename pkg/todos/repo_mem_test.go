package todos

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func TestAddSuccess(t *testing.T) {
	ctx := context.TODO()

	id := uuid.NewString()
	title := "some random title"

	r := NewInMemoryRepository()
	err := r.Add(ctx, Todo{ID: id, Title: title})
	if err != nil {
		t.Fatal(err)
	}

	td, err := r.FindByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if td.Title != title {
		t.Fatal("todo not correctly added")
	}
}

func TestDeleteSuccess(t *testing.T) {
	ctx := context.TODO()

	id := uuid.NewString()
	title := "some random title"

	r := NewInMemoryRepository()
	if err := r.Add(ctx, Todo{ID: id, Title: title}); err != nil {
		t.Fatal(err)
	}

	if err := r.Delete(ctx, id); err != nil {
		t.Fatal(err)
	}

	_, err := r.FindByID(ctx, id)
	if err == nil {
		t.Fatal("todo not deleted")
	}
}

func TestFetchAllAfterConcurrentInserts(t *testing.T) {
	n := rand.Intn(10000)
	ctx := context.TODO()

	r := NewInMemoryRepository()

	g := new(errgroup.Group)
	count := 7
	for idx := 0; idx < count; idx++ {
		g.Go(func() error {
			for i := 0; i < n; i++ {
				if err := r.Add(ctx, Todo{
					ID:    uuid.NewString(),
					Title: fmt.Sprintf("Todo number %d", i+1)}); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	all, err := r.ListAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(all) != n*count {
		t.Fatalf("not all todos were fetched. expected: %d, got: %d", n*count, len(all))
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.TODO()

	id := uuid.NewString()
	title := "Old title"
	newTitle := "New title"

	r := NewInMemoryRepository()
	if err := r.Add(ctx, Todo{ID: id, Title: title}); err != nil {
		t.Fatal(err)
	}

	if err := r.Update(ctx, id, Todo{Title: newTitle}); err != nil {
		t.Fatal(err)
	}

	td, err := r.FindByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if td.Title != newTitle {
		t.Fatal("did not update the todo")
	}
}
