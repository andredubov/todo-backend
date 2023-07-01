package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
)

type todoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem) *todoItemService {
	return &todoItemService{
		repo: repo,
	}
}

func (s *todoItemService) Validate(todolist domain.TodoItem) error {
	return nil
}

func (s *todoItemService) Create(ctx context.Context, listId int, item domain.TodoItem) (int, error) {
	return s.repo.Create(ctx, listId, item)
}

func (s *todoItemService) GetAll(ctx context.Context, userId, listId int) ([]domain.TodoItem, error) {
	return s.repo.GetAll(ctx, userId, listId)
}

func (s *todoItemService) GetById(ctx context.Context, userId, itemId int) (domain.TodoItem, error) {
	return s.repo.GetById(ctx, userId, itemId)
}

func (s *todoItemService) Delete(ctx context.Context, userId, itemId int) error {
	return s.repo.Delete(ctx, userId, itemId)
}

func (s *todoItemService) Update(ctx context.Context, userId, itemId int, input domain.UpdateTodoItemInput) error {
	return s.repo.Update(ctx, userId, itemId, input)
}
