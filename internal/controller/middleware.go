package controller

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/enums"
	"relay-backend/internal/repository"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
	"relay-backend/internal/utils"
	"strconv"
)

const (
	CtxKeyUser CtxKey = iota
)

type CtxKey int8

type AuthMiddleware struct {
	sessionStore *sessions.CookieStore
	userService  *service.UserService
	sessionName  string
	store        *store.Store
}

func ConfigureMiddleware(sessionStore *sessions.CookieStore, sessionName string, store *store.Store) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: sessionStore,
		sessionName:  sessionName,
		store:        store,
	}
}

func (am *AuthMiddleware) Auth(roles []enums.UserRole) func(http.Handler) http.Handler {
	organizationRepository := repository.NewOrganizationRepository(am.store)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := am.sessionStore.Get(r, sessionName)
			if err != nil {
				HTTPError(w, r, utils.NewException(http.StatusInternalServerError, utils.InternalServerError))
				return
			}

			id, ok := s.Values["user_id"]

			if !ok {
				HTTPError(w, r, utils.NewException(http.StatusUnauthorized, utils.Unauthorized))
				return
			}

			if !hasRole(roles, enums.UserRoleEnum.Any) {
				orgId, err := strconv.Atoi(chi.URLParam(r, "orgId"))
				if err != nil {
					HTTPError(w, r, utils.NewException(http.StatusBadRequest, utils.BadRequest))
					return
				}

				userRole := organizationRepository.GetUserRole(id.(int), orgId)

				accessGranted := hasRole(roles, userRole)

				if !accessGranted {
					HTTPError(w, r, utils.NewException(http.StatusForbidden, utils.Forbidden))
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, id)))
		})
	}
}

func hasRole(arr []enums.UserRole, el enums.UserRole) bool {
	for _, val := range arr {
		if val == el {
			return true
		}
	}

	return false
}
