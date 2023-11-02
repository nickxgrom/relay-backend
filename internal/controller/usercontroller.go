package controller

import (
	"encoding/json"
	"net/http"
	"relay-backend/internal/model"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
)

type UserController struct {
	userService *service.UserService
}

type userData struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Patronymic string `json:"patronymic"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

var (
	uc *UserController
)

func UserHandleFunc(s *store.Store) func(w http.ResponseWriter, r *http.Request) {
	if uc == nil {
		uc = &UserController{
			userService: service.NewUserService(s),
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			uc.GetUser(w, r)
		case http.MethodPost:
			uc.CreatUser(w, r)
		}
	}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(CtxKeyUser).(int)
	user, err := uc.userService.FindById(id)
	if err != nil {
		Error(w, r, http.StatusUnauthorized, err)
		return
	}

	Respond(w, r, http.StatusOK, user)
}

func (uc *UserController) CreatUser(w http.ResponseWriter, r *http.Request) {
	ud := &userData{}

	if err := json.NewDecoder(r.Body).Decode(ud); err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

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
}
