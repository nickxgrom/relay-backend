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

func (us *UserService) CreateUser(u *model.User) error {
	err := us.userRepository.Save(u)
	if err != nil {
		return err
	}

	u.Sanitize()

	return nil
}
