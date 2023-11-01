package controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
)

var (
	sessionName = "auth"
)

type SessionController struct {
	sessionService *service.SessionService
	sessionStore   *sessions.CookieStore
}

func NewSessionController(store *store.Store, sessionStore *sessions.CookieStore) *SessionController {
	return &SessionController{
		sessionService: service.NewSessionService(store, sessionStore),
		sessionStore:   sessionStore,
	}
}

func (c SessionController) HandleFunc() func(w http.ResponseWriter, r *http.Request) {
	type session struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &session{}

		//TODO: consider about moving this method to respond.go
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := c.sessionService.CheckUserExist(req.Email, req.Password)
		if err != nil {
			Error(w, r, http.StatusUnauthorized, errors.New("incorrect-email-or-password"))
			return
		}

		s, err := c.sessionStore.Get(r, sessionName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		} else {

		}

		s.Values["user_id"] = u.Id
		err = c.sessionStore.Save(r, w, s)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		Respond(w, r, http.StatusOK, nil)

	}
}