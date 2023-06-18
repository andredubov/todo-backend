package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
	"github.com/andredubov/todo-backend/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	Validate(user domain.User) error
}

type TodoList interface {
	Create(ctx context.Context, todolist domain.TodoList, userId int) (int, error)
	GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (domain.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, todolist domain.TodoList) error
	Validate(list domain.TodoList) error
}

type TodoItem interface {
	Create(ctx context.Context, listId int, item domain.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]domain.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (domain.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, item domain.TodoItem) error
	Validate(item domain.TodoItem) error
}

type Service struct {
	Users
	TodoList
	TodoItem
}

func New(repo *repository.Repository, hasher hash.PasswordHasher) *Service {
	return &Service{
		Users:    NewUsersService(repo.Users, hasher),
		TodoList: NewTodoListService(repo.TodoList),
		TodoItem: NewTodoItemService(repo.TodoItem),
	}
}
