package repository

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Users
	TodoList
	TodoItem
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewPostgresUsersRepository(db),
	}
}
