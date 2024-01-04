package service

import (
	"relay-backend/internal/enums/sender"
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type MessageService struct {
	messageRepository *repository.MessageRepository
}

func NewMessageService(s *store.Store) *MessageService {
	return &MessageService{
		messageRepository: repository.NewMessageRepository(s),
	}
}

func (ms *MessageService) SaveMessage(from sender.Type, text string, senderId int, chatUuid string) error {
	msg := &model.Message{
		Sender:   from,
		Text:     text,
		SenderId: senderId,
		ChatUuid: chatUuid,
	}

	err := ms.messageRepository.Save(msg)
	if err != nil {
		return err
	}

	return nil
}
