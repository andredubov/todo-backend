package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) Validate(domain.User) error {
	return nil
}

func (s *UsersService) Create(ctx context.Context, user domain.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UsersService) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	return s.repo.GetByCredentials(ctx, email, password)
}
