package model

import "time"

type Chat struct {
	Uuid              string    `json:"uuid"`
	WidgetUuid        string    `json:"-"`
	Messages          []Message `json:"messages,omitempty"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
