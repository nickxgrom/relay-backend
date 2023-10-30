package repository

import (
	"relay-backend/internal/model"
	"relay-backend/internal/store"
)

type UserRepository struct {
	store *store.Store
}

func NewUserRepository(s *store.Store) *UserRepository {
	return &UserRepository{
		store: s,
	}
}

func (ur *UserRepository) Save(u *model.User) error {
	err := u.BeforeCreate()
	if err != nil {
		return err
	}

	if err := ur.store.Db.QueryRow(
		"insert into users (first_name, last_name, patronymic, email, encrypted_password) values ($1, $2, $3, $4, $5) returning id",
		&u.FirstName,
		&u.LastName,
		&u.Patronymic,
		&u.Email,
		&u.EncryptedPassword,
	).Scan(&u.Id); err != nil {
		return err
	}

	return nil
}
