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
	"github.com/andredubov/todo-backend/internal/server"
	transport "github.com/andredubov/todo-backend/internal/transport/http"
	"github.com/andredubov/todo-backend/pkg/auth"
	"github.com/andredubov/todo-backend/pkg/logger"
)

const (
	configPath = "./configs"
	timeout    = 5 * time.Second
)

func main() {

	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
		return
	}

	// memcache, hasher := cache.NewMemoryCache(), hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	handler := transport.NewHandler(tokenManager).InitRouts(cfg)

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
