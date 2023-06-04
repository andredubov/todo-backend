package handler

import (
	"net/http"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/gorilla/mux"
)

type Handler struct {
	tokenManager auth.TokenManager
}

func NewHandler(tokenManager auth.TokenManager) *Handler {
	return &Handler{
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRouts(cfg config.Config) http.Handler {

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	return sm
}
