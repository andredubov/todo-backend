package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andredubov/todo-backend/internal/config"
	"github.com/andredubov/todo-backend/internal/repository"
	"github.com/andredubov/todo-backend/internal/server"
	"github.com/andredubov/todo-backend/internal/service"
	transport "github.com/andredubov/todo-backend/internal/transport/http/v1"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/andredubov/todo-backend/pkg/cache"
	"github.com/andredubov/todo-backend/pkg/database"
	"github.com/andredubov/todo-backend/pkg/hash"
	"github.com/andredubov/todo-backend/pkg/logger"
)

const (
	configPath = "./configs"
	timeout    = 5 * time.Second
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Errorf("config initializing failed: %s", err.Error())
		return
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		logger.Errorf("database initializing failed: %s", err.Error())
		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	memcache, hasher := cache.NewMemoryCache(), hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	respository := repository.New(db)
	services := service.New(respository)
	handler := transport.NewHandler(services, tokenManager, hasher, memcache, cfg.Auth.JWT, cfg.CacheTTL).InitRoutes(cfg)

	srv := server.New(cfg, handler)

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
