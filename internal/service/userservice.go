package service

import (
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(s *store.Store) *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(s),
	}
}

func (s *UserService) CreateUser(u *model.User) error {
	err := s.userRepository.Save(u)
	if err != nil {
		return err
	}

	u.Sanitize()

	return nil
}

func (s *UserService) FindByEmail(email string) (*model.User, error) {
	u, err := s.userRepository.Find(email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
