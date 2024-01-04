package repository

import (
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
)

type MessageRepository struct {
	store *store.Store
}

func NewMessageRepository(s *store.Store) *MessageRepository {
	return &MessageRepository{
		store: s,
	}
}

func (mr *MessageRepository) Save(message *model.Message) error {
	err := mr.store.Db.QueryRow(`
		insert into messages (sender, text, chat_uuid, sender_id, timestamp) values ($1::sender, $2, $3, $4, current_timestamp) returning id, timestamp
	`,
		&message.Sender,
		&message.Text,
		&message.ChatUuid,
		&message.SenderId,
	).Scan(&message.Id, &message.Timestamp)

	if err != nil {
		return exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError)
	}

	return nil
}
