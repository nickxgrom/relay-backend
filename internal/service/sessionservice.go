package service

import (
	"errors"
	"relay-backend/internal/model"
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
		userService:       NewUserService(s, nil),
	}
}

func (s *SessionService) CheckUserExist(email string, password string) (*model.User, error) {
	u, err := s.userService.FindByEmail(email)
	if err != nil || !u.ComparePassword(password) {
		return nil, errors.New("")
	}

	return u, nil
}
