package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	Validate(user domain.User) error
}

type TodoList interface {
	Create(ctx context.Context, todolist domain.TodoList, userId int) (int, error)
	GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (domain.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, todolist domain.TodoList) error
	Validate(todolist domain.TodoList) error
}

type TodoItem interface {
	Create(ctx context.Context, listId int, item domain.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]domain.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (domain.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, item domain.TodoItem) error
	Validate(todolist domain.TodoItem) error
}

type Service struct {
	Users
	TodoList
	TodoItem
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Users:    NewUsersService(repo.Users),
		TodoList: NewTodoListService(repo.TodoList),
		TodoItem: NewTodoItemService(repo.TodoItem),
	}
}
