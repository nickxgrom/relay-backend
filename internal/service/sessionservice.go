package service

import (
	"errors"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type SessionService struct {
	sessionRepository *repository.SessionRepository
	userService       *UserService
}

func NewSessionService(s *store.Store) *SessionService {
	return &SessionService{
		sessionRepository: repository.NewSessionRepository(s),
		userService:       NewUserService(s),
	}
}

func (s *SessionService) CreateSession(email string, password string) error {
	u, err := s.userService.FindByEmail(email)
	if err != nil || !u.ComparePassword(password) {
		return errors.New("")
	}

	return nil
}
