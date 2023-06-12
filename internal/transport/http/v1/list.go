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

func (h *Handler) createList(w http.ResponseWriter, r *http.Request) {

	userId := h.getUserId(w, r)

	var todoList domain.TodoList
	if err := json.NewDecoder(r.Body).Decode(&todoList); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	if err := h.services.TodoList.Validate(todoList); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	listId, err := h.services.TodoList.Create(ctx, todoList, userId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable create a todolist"))
		return
	}

	todoList.Id = listId

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(todoList); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

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

func (h *Handler) updateListByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	todoListId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todolist id"))
		return
	}

	var todoList domain.TodoList
	if err := json.NewDecoder(r.Body).Decode(&todoList); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoList.Update(ctx, userId, todoListId, todoList); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to update a todolist by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

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
