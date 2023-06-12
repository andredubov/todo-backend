package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/pkg/errors"
)

const (
	timeout = 5 * time.Second
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body domain.User true "account info"
// @Success 200 {object} domain.User
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {

	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	if err := h.services.Users.Validate(user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid"))
		return
	}

	hash, err := h.passwordHasher.Hash(user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable make password hash"))
		return
	}

	user.Password = hash

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := h.services.Users.Create(ctx, user); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable signup a user"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body domain.User true "credentials"
// @Success 200 {object} SignInResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the given data was not valid JSON"))
		return
	}

	if err := h.services.Users.Validate(user); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, errors.Wrap(err, "the signin request data was not valid"))
		return
	}

	hash, err := h.passwordHasher.Hash(user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "password hashing error"))
		return
	}

	user.Password = hash

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	user, err = h.services.Users.GetByCredentials(ctx, user.Email, user.Password)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable find a user by its credentials"))
		return
	}

	accessToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.Id), h.jwtConfig.AccessTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable create access jwt token"))
		return
	}

	refreshToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.Id), h.jwtConfig.RefreshTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable create refresh jwt token"))
		return
	}

	h.memoryCache.Set(user.Id, user, h.cacheTTL)

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(SignInResponse{AccessToken: accessToken, ResfreshToken: refreshToken}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable encode response data"))
		return
	}
}
