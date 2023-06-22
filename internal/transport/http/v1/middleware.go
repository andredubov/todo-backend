package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/andredubov/todo-backend/internal/domain"
)

const (
	authorizationHeader = "Authorization"
	bearer              = "Bearer"
)

func (h *Handler) getUserId(w http.ResponseWriter, r *http.Request) int {
	user := r.Context().Value(domain.User{}).(domain.User)
	return user.Id
}

func (h *Handler) userIdentity(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		value, err := h.parseAuthHeader(w, r)
		if err != nil {
			h.writeResponseWithError(w, http.StatusUnauthorized, err)
			return
		}

		id, err := strconv.Atoi(value)
		if err != nil {
			h.writeResponseWithError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), domain.User{}, domain.User{Id: id})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) parseAuthHeader(w http.ResponseWriter, r *http.Request) (string, error) {

	header := r.Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != bearer {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}
