package todos

import (
	"context"
	"database/sql"
	"log"
)

type repositoryDB struct {
	conn *sql.DB
}

func NewDbRepository(c *sql.DB) Repository {
	return &repositoryDB{conn: c}
}

func (r *repositoryDB) FindByID(ctx context.Context, id string) (Todo, error) {
	qry := "select * from v_todos where id = ?"
	row := r.conn.QueryRowContext(ctx, qry, id)

	return scan(row)
}

func (r *repositoryDB) ListAll(ctx context.Context) ([]Todo, error) {
	qry := "select * from v_todos"
	rows, err := r.conn.QueryContext(ctx, qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Todo
	for rows.Next() {
		td, err := scan(rows)
		if err != nil {
			log.Printf("Scanning row for ListAll: %v\n", err)
		} else {
			all = append(all, td)
		}
	}

	return all, nil
}

func (r *repositoryDB) Add(ctx context.Context, t Todo) error {
	qry := "insert into todos (id, title) values (?, ?)"

	_, err := r.conn.ExecContext(ctx, qry, t.ID, t.Title)

	return err
}

func (r *repositoryDB) Delete(ctx context.Context, id string) error {
	qry := "delete from todos where id = ?"

	_, err := r.conn.ExecContext(ctx, qry, id)

	return err
}

func (r *repositoryDB) Update(ctx context.Context, id string, t Todo) error {
	qry := "update todos set title = ? where id = ?"

	_, err := r.conn.ExecContext(ctx, qry, t.Title, t.ID)

	return err
}

// Scanner is a constraint that matches sql.Row and sql.Rows
type Scanner interface {
	Scan(...any) error
}

func scan[T Scanner](r T) (Todo, error) {
	var id, title, tags string
	vals := []any{&id, &title, &tags}

	if err := r.Scan(vals...); err != nil {
		return Todo{}, err
	}

	return Todo{
		ID:    id,
		Title: title,
	}, nil
}
