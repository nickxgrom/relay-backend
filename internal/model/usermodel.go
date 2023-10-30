package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int
	FirstName         string
	LastName          string
	Patronymic        string
	Email             string
	Password          string
	EncryptedPassword string
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		encryptedPassword, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = encryptedPassword
	} else {
		return errors.New("password-length-must-be-grater-than-zero")
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
