package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andredubov/todo-backend/internal/domain"
)

const (
	timeout = 5 * time.Second
)

type (
	SignInRequest struct {
		AccessToken   string `json:"email"`
		ResfreshToken string `json:"password"`
	}

	SignInResponse struct {
		AccessToken   string `json:"accessToken"`
		ResfreshToken string `json:"refreshToken"`
	}
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Users.Validate(user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	hash, err := h.passwordHasher.Hash(user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	user.Password = hash

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.Users.Create(ctx, user); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		w.Write([]byte(message))
		return
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Users.Validate(user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	hash, err := h.passwordHasher.Hash(user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	user.Password = hash

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	user, err = h.services.Users.GetByCredentials(ctx, user.Email, user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.ID), h.jwtConfig.AccessTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.ID), h.jwtConfig.RefreshTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	h.memoryCache.Set(user.ID, user, h.cacheTTL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(SignInResponse{AccessToken: accessToken, ResfreshToken: refreshToken}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message := fmt.Sprintf(`{"error": "%s"}`, err.Error())
		w.Write([]byte(message))
		return
	}
}
