package service

import (
	"relay-backend/internal/apiserver/config"
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type UserService struct {
	userRepository *repository.UserRepository
	config         *config.Config
}

func NewUserService(s *store.Store, c *config.Config) *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(s),
		config:         c,
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
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) FindById(id int) (*model.User, error) {
	u, err := s.userRepository.Find(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
