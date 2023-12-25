package model

import "time"

type Message struct {
	Id        int
	From      string //enum "Operator", "Client", "System" ???
	Text      string
	Timestamp time.Time
}
