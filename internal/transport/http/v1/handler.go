package handler

import (
	"fmt"
	"net/http"

	_ "github.com/andredubov/todo-backend/docs"
	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/service"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	jwtConfig    config.JWTConfig
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, jwtConfig config.JWTConfig) *Handler {
	return &Handler{
		tokenManager: tokenManager,
		services:     services,
		jwtConfig:    jwtConfig,
	}
}

func (h *Handler) InitRoutes(cfg config.Config) http.Handler {

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/lists", h.getLists)
	getRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.getListByID)
	getRouter.HandleFunc("/api/lists/{id:[0-9]+}/items", h.getItems)
	getRouter.HandleFunc("/api/items/{id:[0-9]+}", h.getItemByID)
	getRouter.Use(h.userIdentity)

	authRouter := router.Methods(http.MethodPost).Subrouter()
	authRouter.HandleFunc("/auth/sign-up", h.signUp)
	authRouter.HandleFunc("/auth/sign-in", h.signIn)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/api/lists", h.createList)
	postRouter.HandleFunc("/api/lists/{id:[0-9]+}/items", h.createItem)
	postRouter.Use(h.userIdentity)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.updateListByID)
	putRouter.HandleFunc("/api/items/{id:[0-9]+}", h.updateItemByID)
	putRouter.Use(h.userIdentity)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.deleteListByID)
	deleteRouter.HandleFunc("/api/items/{id:[0-9]+}", h.deleteItemByID)
	deleteRouter.Use(h.userIdentity)

	return router
}

func (h *Handler) writeResponseWithError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	message := fmt.Sprintf(`{"message": "%s"}`, err.Error())
	w.Write([]byte(message))
}

func (h *Handler) writeResponseHeader(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
