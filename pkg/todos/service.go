package todos

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	FindByID(context.Context, string) (Todo, error)
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
