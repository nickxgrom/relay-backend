package repository

import (
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
)

type ChatRepository struct {
	store *store.Store
}

func NewChatRepository(s *store.Store) *ChatRepository {
	return &ChatRepository{
		store: s,
	}
}

func (cr *ChatRepository) Save(chat *model.Chat) error {
	err := cr.store.Db.QueryRow(`insert into chats (uuid, creation_timestamp, widget_uuid) values ($1, current_timestamp, $2) returning creation_timestamp`,
		&chat.Uuid,
		&chat.WidgetUuid,
	).Scan(&chat.CreationTimestamp)

	if err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	return nil
}
