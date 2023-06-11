package handler

import (
	"fmt"
	"net/http"

	"github.com/andredubov/todo-backend/internal/domain"
)

func (h *Handler) createList(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.User{}).(domain.User)
	message := fmt.Sprintf("createList userID=%d", user.ID)
	w.Write([]byte(message))
}

func (h *Handler) getAllList(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.User{}).(domain.User)
	message := fmt.Sprintf("getAllList userID=%d", user.ID)
	w.Write([]byte(message))
}

func (h *Handler) getListByID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.User{}).(domain.User)
	message := fmt.Sprintf("getListByID userID=%d", user.ID)
	w.Write([]byte(message))
}

func (h *Handler) updateListByID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.User{}).(domain.User)
	message := fmt.Sprintf("updateListByID userID=%d", user.ID)
	w.Write([]byte(message))
}

func (h *Handler) deleteListByID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(domain.User{}).(domain.User)
	message := fmt.Sprintf("deleteListByID userID=%d", user.ID)
	w.Write([]byte(message))
}
