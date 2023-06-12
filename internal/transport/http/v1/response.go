package handler

import "github.com/andredubov/todo-backend/internal/domain"

const (
	success = "success"
	fail    = "fail"
)

type (
	GetTodoListsResponse struct {
		Data []domain.TodoList `json:"data"`
	}

	SignInResponse struct {
		AccessToken   string `json:"accessToken"`
		ResfreshToken string `json:"refreshToken"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
	}

	StatusResponse struct {
		Status string `json:"status"`
	}
)
