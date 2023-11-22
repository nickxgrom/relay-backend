package model

import "relay-backend/internal/enums"

type Employee struct {
	Id       int            `json:"id"`
	UserRole enums.UserRole `json:"userRole"`
}
