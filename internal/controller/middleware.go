package controller

import (
	"context"
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/service"
)

const (
	CtxKeyUser CtxKey = iota
)

type CtxKey int8

func AuthenticateUser(sessionStore *sessions.CookieStore, userService *service.UserService, sessionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := sessionStore.Get(r, sessionName)
			if err != nil {
				Error(w, r, http.StatusInternalServerError, err)
				return
			}

			id, ok := s.Values["user_id"]

			if !ok {
				Error(w, r, http.StatusUnauthorized, errors.New("not-authenticated"))
				return
			}

			u, err := userService.FindById(id.(int))
			if err != nil {
				Error(w, r, http.StatusUnauthorized, errors.New("not-authenticated"))
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, u)))
		})
	}
}
