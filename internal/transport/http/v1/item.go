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

// @Summary Create todo item
// @Security ApiKeyAuth
// @Tags items
// @Description create todo item
// @ID create-item
// @Accept json
// @Produce json
// @Param input body domain.TodoItem true "list info"
// @Success 200 {object} domain.TodoItem
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists [post]
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
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	itemId, err := h.services.TodoItem.Create(ctx, listId, todoItem)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to create a todo-item"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(domain.TodoItem{Id: itemId}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary Get All Items
// @Security ApiKeyAuth
// @Tags items
// @Description get all todo-items
// @ID get-all-items
// @Accept json
// @Produce json
// @Success 200 {object} GetTodoItemResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/lists/:id/items [get]
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
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to find any todolists by user id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(GetTodoItemResponse{Data: todoItems}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary Get todo-item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description get todo-item by id
// @ID get-item-by-id
// @Accept json
// @Produce json
// @Success 200 {object} domain.TodoItem
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/items/:id [get]
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
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(todoItem); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary Update todo-item by Id
// @Security ApiKeyAuth
// @Tags items
// @Description update todo-item by id
// @ID update-item-by-id
// @Accept json
// @Produce json
// @Param input body domain.UpdateTodoItemInput true "item info"
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/items/:id [put]
func (h *Handler) updateItemByID(w http.ResponseWriter, r *http.Request) {

	userId, vars := h.getUserId(w, r), mux.Vars(r)

	itemId, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to convert a todo-item id"))
		return
	}

	var updateTodoItemInput domain.UpdateTodoItemInput
	if err := json.NewDecoder(r.Body).Decode(&updateTodoItemInput); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.TodoItem.Update(ctx, userId, itemId, updateTodoItemInput); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to update a todo-item by id"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(StatusResponse{success}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary Delete todo-item by Id
// @Security ApiKeyAuth
// @Tags items
// @Description delete todo-item by id
// @ID delete-item-by-id
// @Accept json
// @Produce json
// @Success 200 {object} StatusResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/items/:id [delete]
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
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}
