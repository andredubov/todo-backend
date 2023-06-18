package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
	"github.com/andredubov/todo-backend/pkg/hash"
)

type UsersService struct {
	repo           repository.Users
	passwordHasher hash.PasswordHasher
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repo:           repo,
		passwordHasher: hasher,
	}
}

func (s *UsersService) Validate(user domain.User) error {
	return nil
}

func (s *UsersService) Create(ctx context.Context, user domain.User) (int, error) {

	hash, err := s.passwordHasher.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	user.Password = hash

	return s.repo.Create(ctx, user)
}

func (s *UsersService) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {

	hash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}

	return s.repo.GetByCredentials(ctx, email, hash)
}
