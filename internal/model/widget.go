package model

type Widget struct {
	Organization *Organization `json:"-"`
	Uuid         string        `json:"uuid"`
}
