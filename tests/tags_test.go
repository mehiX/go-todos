package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mehix/go-todos/pkg/todos"
)

func TestSearchByTags(t *testing.T) {

	ctx := context.Background()

	tags := []string{
		uuid.NewString(),
		uuid.NewString(),
		uuid.NewString(),
	}

	addTodo(ctx, todos.Todo{ID: uuid.NewString(), Title: "some title", Tags: []string{tags[0]}})
	addTodo(ctx, todos.Todo{ID: uuid.NewString(), Title: "some title 2", Tags: []string{tags[0], tags[1]}})
	addTodo(ctx, todos.Todo{ID: uuid.NewString(), Title: "some title 3", Tags: []string{tags[2]}})

	withTag1, err := byTags([]string{tags[0]})
	if err != nil {
		t.Fatal(err)
	}

	if withTag1 != 2 {
		t.Fatalf("wrong number of results. expected: %d, got: %d", 2, withTag1)
	}

	withTag3, err := byTags([]string{tags[2]})
	if err != nil {
		t.Fatal(err)
	}

	if withTag3 != 1 {
		t.Fatalf("wrong number of results. expected: %d, got: %d", 1, withTag3)
	}

	withAllTags, err := byTags(tags)
	if err != nil {
		t.Fatal(err)
	}
	if withAllTags != 3 {
		t.Fatalf("wrong number of results. expected: %d, got: %d", 3, withAllTags)
	}
}
