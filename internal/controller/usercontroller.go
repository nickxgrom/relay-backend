package controller

import (
	"io"
	"net/http"
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

func (uc *UserController) HandleFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uc.CreateUser(w, r)
	default:
		io.WriteString(w, "method-not-allowed")
		w.WriteHeader(400)
	}

}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := uc.userService.CreateUser()
	if err != nil {
		io.WriteString(w, err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
