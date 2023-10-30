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

func (us *UserService) CreateUser() error {
	err := us.userRepository.Save(&model.User{
		FirstName: "anon",
		LastName:  "anon",
		Email:     "super@puper.org",
		Password:  "5658",
	})
	if err != nil {
		return err
	}

	return nil
}
