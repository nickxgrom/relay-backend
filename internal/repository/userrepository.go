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
		"insert into users (first_name, last_name, patronymic, email, encrypted_password) values ($1, $2, $3, $4, $5) returning id, verified",
		&u.FirstName,
		&u.LastName,
		&u.Patronymic,
		&u.Email,
		&u.EncryptedPassword,
	).Scan(&u.Id, &u.Verified); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user-not-found")
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}

	if err := ur.store.Db.QueryRow(
		"select id, first_name, last_name, patronymic, email, encrypted_password from users where id = $1",
		id,
	).Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Patronymic,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user-not-found")
		}

		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) SaveToken(userId int, token string) {
	ur.store.Db.QueryRow("insert into email_tokens (user_id, token) values ($1, $2)", userId, token)
}
