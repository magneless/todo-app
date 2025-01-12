package reposiory

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Auth interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Reposiory struct {
	Auth
	TodoList
	TodoItem
}

func New(db *sql.DB) *Reposiory {
	return &Reposiory{}
}
