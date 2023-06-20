package todos

import "fmt"

type ErrNotFound struct {
	id string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found todo with id: %s", e.id)
}
