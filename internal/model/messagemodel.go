package model

import (
	"relay-backend/internal/enums/sender"
	"time"
)

type Message struct {
	Id        int         `json:"id"`
	ChatUuid  string      `json:"chatUuid"`
	Sender    sender.Type `json:"sender"`
	SenderId  int         `json:"senderId,omitempty"`
	Text      string      `json:"text"`
	Timestamp time.Time   `json:"timestamp"`
}
