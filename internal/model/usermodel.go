package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int    `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Patronymic        string `json:"patronymic,omitempty"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Verified          bool   `json:"verified"`
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

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (u *User) Sanitize() {
	u.Password = ""
}
