package model

import "time"

type Chat struct {
	Uuid              string    `json:"uuid"`
	WidgetUuid        string    `json:"-"`
	Messages          []Message `json:"messages,omitempty"`
	Archived          bool      `json:"archived"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
