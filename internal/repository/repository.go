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
	Create(ctx context.Context, todolist domain.TodoList, userId int) error
	GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (domain.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input domain.TodoList) error
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
		Users:    NewPostgresUsersRepository(db),
		TodoList: NewPostgresTodoListRepository(db),
	}
}
