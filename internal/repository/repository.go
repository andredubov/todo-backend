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
	Create(ctx context.Context, todolist domain.TodoList, userId int) (int, error)
	GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (domain.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input domain.TodoList) error
}

type TodoItem interface {
	Create(ctx context.Context, listId int, item domain.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]domain.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (domain.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, item domain.TodoItem) error
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
		TodoItem: NewPostgresTodoItemRepository(db),
	}
}
