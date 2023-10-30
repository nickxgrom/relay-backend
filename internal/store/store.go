package store

import "database/sql"

type Store struct {
	Db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		Db: db,
	}
}
