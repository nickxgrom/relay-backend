package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
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

func NewUserController(s *store.Store, sessionStore *sessions.CookieStore) func(r chi.Router) {
	if uc == nil {
		uc = &UserController{
			userService: service.NewUserService(s),
		}
	}

	am := ConfigureMiddleware(sessionStore, "auth")

	return func(r chi.Router) {
		r.Post("/", uc.CreateUser)
		r.With(am.AuthenticateUser).Get("/", uc.GetUser)
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

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
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
