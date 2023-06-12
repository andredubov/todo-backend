package service

import (
	"context"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/andredubov/todo-backend/internal/repository"
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
	return nil
}

func (s *todoListService) Create(ctx context.Context, todolist domain.TodoList, userID int) error {
	return s.repo.Create(ctx, todolist, userID)
}

func (s *todoListService) GetByUserId(ctx context.Context, userId int) ([]domain.TodoList, error) {
	return s.repo.GetByUserId(ctx, userId)
}

func (s *todoListService) GetById(ctx context.Context, userId, listId int) (domain.TodoList, error) {
	return s.repo.GetById(ctx, userId, listId)
}

func (s *todoListService) Update(ctx context.Context, userId, listId int, todolist domain.TodoList) error {
	return s.repo.Update(ctx, userId, listId, todolist)
}

func (s *todoListService) Delete(ctx context.Context, userId, listId int) error {
	return s.repo.Delete(ctx, userId, listId)
}
