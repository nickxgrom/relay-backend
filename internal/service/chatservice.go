package service

import (
	"github.com/google/uuid"
	"relay-backend/internal/model"
	"relay-backend/internal/repository"
	"relay-backend/internal/store"
)

type ChatService struct {
	chatRepository *repository.ChatRepository
}

func NewChatService(s *store.Store) *ChatService {
	return &ChatService{
		chatRepository: repository.NewChatRepository(s),
	}
}

func (cs *ChatService) CreateNewChat(widgetUuid string) (*model.Chat, error) {
	chat := &model.Chat{
		Uuid:       uuid.NewString(),
		WidgetUuid: widgetUuid,
	}

	err := cs.chatRepository.Save(chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
