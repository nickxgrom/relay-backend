package repository

import (
	"database/sql"
	"errors"
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

func (ur *UserRepository) Find(email string) (*model.User, error) {
	u := &model.User{}

	if err := ur.store.Db.QueryRow(
		"select id, first_name, last_name, patronymic, email, encrypted_password from users where email = $1",
		email,
	).Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Patronymic,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user-not-found")
		}

		return nil, err
	}

	return u, nil
}
