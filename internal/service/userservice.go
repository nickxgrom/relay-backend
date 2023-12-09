package service

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/smtp"
	"regexp"
	"relay-backend/internal/apiserver/config"
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
)

type UserService struct {
	userRepository *repository.UserRepository
	config         *config.Config
}

const (
	emailRegex = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
)

func NewUserService(s *store.Store, c *config.Config) *UserService {
	return &UserService{
		userRepository: repository.NewUserRepository(s),
		config:         c,
	}
}

func (s *UserService) CreateUser(u *model.User) error {
	if ok, _ := regexp.Match(emailRegex, []byte(u.Email)); ok {
		if err := s.userRepository.Save(u); err != nil {
			return err
		}
	} else {
		return exception.NewException(http.StatusBadRequest, exception.Enum.InvalidEmail)
	}

	token := uuid.NewString()
	msg := fmt.Sprintf("Subject: Relay email confirmation token\n\rRelay confirmation system introduces email confirmation token:\n\r%s", token)

	s.userRepository.SaveToken(u.Id, token)

	if err := s.SendEmail(u.Email, msg); err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	u.Sanitize()

	return nil
}

func (s *UserService) UpdateUser(userId int, user *model.User) error {
	if user.Email != "" {
		ok, _ := regexp.Match(emailRegex, []byte(user.Email))
		if !ok {
			return exception.NewException(http.StatusBadRequest, exception.Enum.InvalidEmail)
		}
	}

	u, err := s.FindById(userId)
	if err != nil {
		return err
	}
	newEmail := user.Email

	if err := s.userRepository.Update(userId, user); err != nil {
		return err
	}

	if user.Email != "" && newEmail != u.Email {
		err := s.userRepository.SetVerified(userId, false)
		if err != nil {
			return err
		}
		user.Verified = false

		if err := s.userRepository.DeleteAllTokens(userId); err != nil {
			return err
		}
		//TODO: new token for email
		//TODO: notification for old email, that email was changed
	}

	defer user.Sanitize()
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

func (s *UserService) SendEmail(email string, message string) error {
	smtpEmail := s.config.SmtpEmail
	smtpPassword := s.config.SmtpPassword

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	smtpAuth := smtp.PlainAuth("", smtpEmail, smtpPassword, smtpHost)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpHost, smtpPort),
		smtpAuth,
		smtpEmail,
		[]string{email},
		[]byte(message),
	)
}

func (s *UserService) ConfirmEmail(userId int, token string) error {
	if err := s.userRepository.FindToken(userId, token); err != nil {
		return err
	}

	if err := s.userRepository.SetVerified(userId, true); err != nil {
		return err
	}

	if err := s.userRepository.DeleteToken(userId, token); err != nil {
		return err
	}

	return nil
}
