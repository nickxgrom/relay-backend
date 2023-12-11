package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
)

var (
	sessionName  = "auth"
	c            *SessionController
	cookieMaxAge = 60 * 60 * 24 // day
)

type SessionController struct {
	sessionService *service.SessionService
	sessionStore   *sessions.CookieStore
}

func NewSessionController(s *store.Store, sessionStore *sessions.CookieStore) func(r chi.Router) {
	if c == nil {
		c = &SessionController{
			sessionService: service.NewSessionService(s),
			sessionStore:   sessionStore,
		}
	}

	type session struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(router chi.Router) {
		router.Post("/", func(w http.ResponseWriter, r *http.Request) {
			req := &session{}

			//TODO: consider about moving this method to respond.go
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				HTTPError(w, r, exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError))
				return
			}

			u, err := c.sessionService.CheckUserExist(req.Email, req.Password)
			if err != nil {
				HTTPError(w, r, exception.NewException(http.StatusUnauthorized, exception.Enum.IncorrectEmailOrPassword))
				return
			}

			s, err := c.sessionStore.Get(r, sessionName)
			if err != nil {
				HTTPError(w, r, exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError))
				return
			} else {

			}

			s.Values["user_id"] = u.Id
			s.Options.MaxAge = cookieMaxAge
			err = c.sessionStore.Save(r, w, s)
			if err != nil {
				HTTPError(w, r, exception.NewException(http.StatusUnauthorized, exception.Enum.IncorrectEmailOrPassword))
				return
			}

			Respond(w, r, http.StatusOK, nil)
		})
	}
}
