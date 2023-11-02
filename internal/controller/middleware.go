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

type AuthMiddleware struct {
	sessionStore *sessions.CookieStore
	userService  *service.UserService
	sessionName  string
}

func ConfigureMiddleware(sessionStore *sessions.CookieStore, sessionName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: sessionStore,
		sessionName:  sessionName,
	}
}

func (am *AuthMiddleware) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := am.sessionStore.Get(r, sessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := s.Values["user_id"]

		if !ok {
			Error(w, r, http.StatusUnauthorized, errors.New("not-authenticated"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, id)))
	})
}
