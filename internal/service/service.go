package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	// GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Validate(domain.User) error
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Users
	TodoList
	TodoItem
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Users: NewUsersService(repo.Users),
	}
}
