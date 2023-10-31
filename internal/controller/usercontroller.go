package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(s *store.Store) *UserController {
	return &UserController{
		userService: service.NewUserService(s),
	}
}

func (uc *UserController) HandleFunc() func(w http.ResponseWriter, r *http.Request) {
	type userData struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Patronymic string `json:"patronymic"`
		Email      string `json:"email"`
		Password   string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ud := &userData{}

		if err := json.NewDecoder(r.Body).Decode(ud); err != nil {
			Error(w, r, http.StatusBadRequest, err)
			return
		}

		switch r.Method {
		case "POST":
			u := &model.User{
				FirstName:  ud.FirstName,
				LastName:   ud.LastName,
				Patronymic: ud.Patronymic,
				Email:      ud.Email,
				Password:   ud.Password,
			}

			err := uc.userService.CreateUser(u)

			if err != nil {
				Error(w, r, http.StatusBadRequest, err)
				return
			}

			Respond(w, r, http.StatusCreated, u)
		default:
			Error(w, r, http.StatusMethodNotAllowed, errors.New("method-not-allowed"))
		}
	}
}
