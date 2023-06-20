package todos

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	FindByID(context.Context, string) (Todo, error)
	FindByTags(context.Context, []string) ([]Todo, error)
	ListAll(context.Context) ([]Todo, error)
	Add(context.Context, Todo) (Todo, error)
	Delete(context.Context, string) error
	Update(context.Context, string, Todo) (Todo, error)
}

type service struct {
	repo Repository
}

func NewService(opts ...Option) Service {
	s := &service{}
	for _, o := range opts {
		o(s)
	}
	return s
}

type Option func(*service)

func WithRepo(r Repository) Option {
	return func(s *service) {
		s.repo = r
	}
}

func (s *service) FindByID(ctx context.Context, id string) (Todo, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListAll(ctx context.Context) ([]Todo, error) {
	return s.repo.ListAll(ctx)
}

func (s *service) Add(ctx context.Context, t Todo) (Todo, error) {
	t.ID = uuid.NewString()
	if err := s.repo.Add(ctx, t); err != nil {
		return Todo{}, err
	}

	return s.FindByID(ctx, t.ID)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("provided ID is not a UUID")
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, t Todo) (Todo, error) {
	if _, err := uuid.Parse(id); err != nil {
		return Todo{}, fmt.Errorf("provided ID is not a UUID")
	}

	if err := s.repo.Update(ctx, id, t); err != nil {
		return Todo{}, err
	}

	return s.FindByID(ctx, id)
}

// FindByTags returns all the todo's that contain at least one of the provided flags
func (s *service) FindByTags(ctx context.Context, tags []string) ([]Todo, error) {

	ch := make(chan Todo)
	errs := make(chan error, 1)
	unique := make(map[string]Todo)

	g, gCtx := errgroup.WithContext(ctx)

	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		// https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func(t string) func() error {
			return func() error {
				tds, err := s.repo.FindByTag(ctx, t)
				if err != nil {
					return err
				}
				for _, td := range tds {
					select {
					case <-gCtx.Done():
						return gCtx.Err()
					case ch <- td:
					}
				}

				return nil
			}
		}(t))
	}

	go func() {
		defer close(ch)
		if err := g.Wait(); err != nil {
			log.Printf("Finding by tags: %v", err)
			errs <- err
		}
		close(errs)
	}()

	for t := range ch {
		unique[t.ID] = t
	}

	var all []Todo
	for _, v := range unique {
		all = append(all, v)
	}

	return all, <-errs
}
