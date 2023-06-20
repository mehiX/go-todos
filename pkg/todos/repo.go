package todos

import "context"

type Repository interface {
	FindByID(context.Context, string) (Todo, error)
	ListAll(context.Context) ([]Todo, error)
	Add(context.Context, Todo) error
	Delete(context.Context, string) error
	Update(context.Context, string, Todo) error
}
