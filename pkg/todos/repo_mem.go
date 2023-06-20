package todos

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/exp/slices"
)

type repositoryMem struct {
	data map[string]Todo
	m    sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &repositoryMem{
		data: make(map[string]Todo),
	}
}

func (r *repositoryMem) Add(_ context.Context, td Todo) error {

	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.data[td.ID]; ok {
		return fmt.Errorf("a todo with this ID already exists")
	}
	r.data[td.ID] = td

	return nil
}

func (r *repositoryMem) Delete(_ context.Context, id string) error {
	r.m.Lock()
	defer r.m.Unlock()

	delete(r.data, id)

	return nil
}

func (r *repositoryMem) Update(_ context.Context, id string, td Todo) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("a todo with this ID does not exist")
	}

	r.data[id] = td

	return nil
}

func (r *repositoryMem) ListAll(_ context.Context) ([]Todo, error) {

	all := make([]Todo, 0)

	r.m.RLock()
	for _, td := range r.data {
		all = append(all, td)
	}
	r.m.RUnlock()

	return all, nil
}

func (r *repositoryMem) FindByID(_ context.Context, id string) (Todo, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	td, ok := r.data[id]
	if !ok {
		return Todo{}, fmt.Errorf("a todo with this ID does not exist")
	}

	return td, nil
}

func (r *repositoryMem) FindByTag(_ context.Context, tg string) ([]Todo, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	all := make([]Todo, 0)

	for _, v := range r.data {
		if slices.Contains(v.Tags, tg) {
			all = append(all, v)
		}
	}

	return all, nil
}
