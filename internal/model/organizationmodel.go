package model

type Organization struct {
	Id           int    `json:"id"`
	OwnerId      int    `json:"ownerId"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Widget       string `json:"widget"`
	CreationDate string `json:"creationDate"`
	Employees    []User `json:"employees,omitempty"`
}
