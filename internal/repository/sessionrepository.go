package repository

import "relay-backend/internal/store"

type SessionRepository struct {
	store *store.Store
}

func NewSessionRepository(s *store.Store) *SessionRepository {
	return &SessionRepository{
		store: s,
	}
}
