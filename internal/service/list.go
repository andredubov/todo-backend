package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
	"gopkg.in/validator.v2"
)

type todoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *todoListService {
	return &todoListService{
		repo: repo,
	}
}

func (s *todoListService) Validate(todolist domain.TodoList) error {

	if err := validator.Validate(todolist); err != nil {
		return err
	}

	return nil
}

func (s *todoListService) Create(ctx context.Context, todolist domain.TodoList, userID int) (int, error) {
	return s.repo.Create(ctx, todolist, userID)
}

func (s *todoListService) GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error) {
	return s.repo.GetByUserId(ctx, userId)
}

func (s *todoListService) GetById(ctx context.Context, userId, listId int) (domain.TodoList, error) {
	return s.repo.GetById(ctx, userId, listId)
}

func (s *todoListService) Update(ctx context.Context, userId, listId int, input domain.UpdateTodoListInput) error {
	return s.repo.Update(ctx, userId, listId, input)
}

func (s *todoListService) Delete(ctx context.Context, userId, listId int) error {
	return s.repo.Delete(ctx, userId, listId)
}
