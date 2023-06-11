package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/service"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/andredubov/todo-backend/pkg/cache"
	"github.com/andredubov/todo-backend/pkg/hash"
	"github.com/gorilla/mux"
)

type Handler struct {
	services       *service.Service
	tokenManager   auth.TokenManager
	passwordHasher hash.PasswordHasher
	memoryCache    cache.Cache
	jwtConfig      config.JWTConfig
	cacheTTL       time.Duration
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, passwordHasher hash.PasswordHasher, memoryCache cache.Cache, jwtConfig config.JWTConfig, cacheTTL time.Duration) *Handler {
	return &Handler{
		tokenManager:   tokenManager,
		services:       services,
		passwordHasher: passwordHasher,
		memoryCache:    memoryCache,
		jwtConfig:      jwtConfig,
		cacheTTL:       cacheTTL,
	}
}

func (h *Handler) InitRoutes(cfg config.Config) http.Handler {

	router := mux.NewRouter()

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/lists", h.getAllList)
	getRouter.HandleFunc("/api/lists/{id:[0-9]+}", h.getListByID)
	getRouter.HandleFunc("/api/lists/{id:[0-9]+}/items", h.getAllItems)
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
	message := fmt.Sprintf(`{"error": "%s"}`, err.Error())
	w.Write([]byte(message))
}
