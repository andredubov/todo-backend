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

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	listId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-list id"))
		return
	}

	var todoItem domain.TodoItem
	if err := json.NewDecoder(r.Body).Decode(&todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	if err := h.services.TodoItem.Validate(todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	itemId, err := h.services.TodoItem.Create(ctx, listId, todoItem)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable create a todo-item"))
		return
	}

	todoItem.Id = itemId

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

func (h *Handler) getItems(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	listId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-list id"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	todoItems, err := h.services.TodoItem.GetAll(ctx, userId, listId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable find any todolists by user id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(GetTodoItemResponse{Data: todoItems}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

func (h *Handler) getItemByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	itemId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-item id"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	todoItem, err := h.services.TodoItem.GetById(ctx, userId, itemId)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable find any todo-item by user id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

func (h *Handler) updateItemByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	itemId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-item id"))
		return
	}

	var todoItem domain.TodoItem
	if err := json.NewDecoder(r.Body).Decode(&todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoItem.Update(ctx, userId, itemId, todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to update a todo-item by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

func (h *Handler) deleteItemByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	itemId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-item id"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoItem.Delete(ctx, userId, itemId); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to delete a todo-item by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}
