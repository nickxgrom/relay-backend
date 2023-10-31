package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
)

type SessionController struct {
	sessionService *service.SessionService
}

func NewSessionController(s *store.Store) *SessionController {
	return &SessionController{
		sessionService: service.NewSessionService(s),
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

		err := c.sessionService.CreateSession(req.Email, req.Password)

		if err != nil {
			Error(w, r, http.StatusUnauthorized, errors.New("incorrect-email-or-password"))
			return
		}

		Respond(w, r, http.StatusOK, nil)

	}
}
