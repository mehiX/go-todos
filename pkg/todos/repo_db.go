package todos

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	qry := "insert into todos (id, title, tags) values (?, ?, ?)"

	_, err := r.conn.ExecContext(ctx, qry, t.ID, t.Title, t.CleanTags())

	return err
}

func (r *repositoryDB) Delete(ctx context.Context, id string) error {
	qry := "delete from todos where id = ?"

	_, err := r.conn.ExecContext(ctx, qry, id)

	return err
}

func (r *repositoryDB) Update(ctx context.Context, id string, t Todo) error {
	fmt.Printf("Updated todo: %#v\n", t)
	qry := "update todos set title = ?, tags = ?, completed_at = ? where id = ?"

	_, err := r.conn.ExecContext(ctx, qry, t.Title, t.CleanTags(), sql.NullTime{Time: *t.CompletedAt, Valid: t.CompletedAt != nil}, t.ID)
	if err != nil {
		fmt.Printf("Update error: %v", err)
	}

	return err
}

// FindByTag returns all Todo's that contain this tag
// Logic should be improved, maybe refactoring the whole tags logic. At the moment it will return
// entries containing `golang` when `go` is passed in as parameter
func (r *repositoryDB) FindByTag(ctx context.Context, tg string) ([]Todo, error) {
	qry := "select * from v_todos where tags like ?"

	rows, err := r.conn.QueryContext(ctx, qry, "%"+tg+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Todo
	for rows.Next() {
		td, err := scan(rows)
		if err != nil {
			return all, err
		} else {
			all = append(all, td)
		}
	}

	return all, nil

}

// Scanner is a constraint that matches sql.Row and sql.Rows
type Scanner interface {
	Scan(...any) error
}

func scan[T Scanner](r T) (Todo, error) {
	var id, title, tags string
	var completedAt sql.NullTime
	vals := []any{&id, &title, &tags, &completedAt}

	if err := r.Scan(vals...); err != nil {
		log.Printf("Scan error: %v\n", err)
		return Todo{}, err
	}

	completedWhen := &completedAt.Time
	if !completedAt.Valid {
		completedWhen = nil
	}

	return Todo{
		ID:          id,
		Title:       title,
		Tags:        strings.Split(tags, ","),
		CompletedAt: completedWhen,
	}, nil
}
