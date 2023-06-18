package service

import (
	"context"
	"net/mail"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
	"github.com/andredubov/todo-backend/pkg/hash"
	"github.com/pkg/errors"
	"gopkg.in/validator.v2"
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

	if err := validator.Validate(user); err != nil {
		return errors.Wrap(err, "sign-up input data not valid")
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return errors.Wrap(err, "sign-up input data not valid")
	}

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
