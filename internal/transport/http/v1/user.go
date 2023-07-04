package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/andredubov/todo-backend/internal/domain"
	"github.com/pkg/errors"
	"gopkg.in/validator.v2"
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
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	userId, err := h.services.Users.Create(ctx, user)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to signup a user"))
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(domain.User{Id: userId}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, errors.Wrap(err, "unable to encode response data"))
		return
	}
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body domain.Credentials true "credentials"
// @Success 200 {object} SignInResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {

	var credentials domain.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.Validate(credentials); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	user, err := h.services.Users.GetByCredentials(ctx, credentials)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	accessToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.Id), h.jwtConfig.AccessTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := h.tokenManager.NewJWT(strconv.Itoa(user.Id), h.jwtConfig.RefreshTokenTTL)
	if err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	h.writeResponseHeader(w, http.StatusOK)

	if err := json.NewEncoder(w).Encode(SignInResponse{AccessToken: accessToken, ResfreshToken: refreshToken}); err != nil {
		h.writeResponseWithError(w, http.StatusInternalServerError, err)
		return
	}
}
