package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body domain.TodoList true "list info"
// @Success 200 {object} domain.TodoList
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists [post]
func (h *Handler) createList(w http.ResponseWriter, r *http.Request) {

	userId := h.getUserId(w, r)

	var todoList domain.TodoList
	if err := json.NewDecoder(r.Body).Decode(&todoList); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	if err := h.services.TodoList.Validate(todoList); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	listId, err := h.services.TodoList.Create(ctx, todoList, userId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to create a todolist"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(domain.TodoList{Id: listId}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all todo-lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} GetTodoListsResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists [get]
func (h *Handler) getLists(w http.ResponseWriter, r *http.Request) {

	userId := h.getUserId(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	todolists, err := h.services.TodoList.GetByUserId(ctx, userId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable find any todolists by user id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(GetTodoListsResponse{Data: todolists}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get todo-list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.TodoList
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists/:id [get]
func (h *Handler) getListByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	todoListId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todolist id"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	todoList, err := h.services.TodoList.GetById(ctx, userId, todoListId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to get a todolist by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

// @Summary Update todo-list by Id
// @Security ApiKeyAuth
// @Tags lists
// @Description update todo-list by id
// @ID update-list-by-id
// @Accept json
// @Produce json
// @Param input body domain.UpdateTodoListInput true "item info"
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists/:id [put]
func (h *Handler) updateListByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	todoListId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todolist id"))
		return
	}

	var updateTodoListInput domain.UpdateTodoListInput
	if err := json.NewDecoder(r.Body).Decode(&updateTodoListInput); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoList.Update(ctx, userId, todoListId, updateTodoListInput); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to update a todolist by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

// @Summary Delete todo-list by Id
// @Security ApiKeyAuth
// @Tags lists
// @Description delete todo-list by id
// @ID delete-list-by-id
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists/:id [delete]
func (h *Handler) deleteListByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	todoListId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todolist id"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoList.Delete(ctx, userId, todoListId); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to delete a todolist by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}
