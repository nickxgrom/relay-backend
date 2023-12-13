package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"relay-backend/internal/apiserver/config"
	"relay-backend/internal/enums"
	"relay-backend/internal/model"
	"relay-backend/internal/service"
	"relay-backend/internal/store"
	"relay-backend/internal/utils/exception"
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

func NewUserController(store *store.Store, middleware *AuthMiddleware, config *config.Config) func(r chi.Router) {
	if uc == nil {
		uc = &UserController{
			userService: service.NewUserService(store, config),
		}
	}

	return func(r chi.Router) {
		r.Post("/", uc.CreateUser)
		r.Post("/forgot-password", uc.forgotPassword)
		r.With(middleware.Auth(enums.Access.Any)).Get("/", uc.GetUser)
		r.With(middleware.Auth(enums.Access.Any)).Put("/", uc.UpdateUser)
		r.With(middleware.Auth(enums.Access.Any)).Post("/confirm-email", uc.ConfirmEmail)
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

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(CtxKeyUser).(int)

	ud := &userData{}
	if err := json.NewDecoder(r.Body).Decode(ud); err != nil {
		HTTPError(w, r, exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError))
		return
	}

	user := model.User{
		FirstName:  ud.FirstName,
		LastName:   ud.LastName,
		Patronymic: ud.Patronymic,
		Email:      ud.Email,
		Password:   ud.Password,
	}

	err := uc.userService.UpdateUser(id, &user)
	if err != nil {
		HTTPError(w, r, err.(exception.Exception))
		return
	}

	Respond(w, r, http.StatusOK, user)
}

func (uc *UserController) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(CtxKeyUser).(int)
	type data struct {
		Token string `json:"token"`
	}

	body := &data{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		HTTPError(w, r, exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError))
		return
	}

	err := uc.userService.ConfirmEmail(id, body.Token)
	if err != nil {
		HTTPError(w, r, err.(exception.Exception))
		return
	}

	Respond(w, r, http.StatusOK, nil)
}

func (uc *UserController) forgotPassword(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Email string `json:"email"`
	}

	body := &data{}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		HTTPError(w, r, exception.NewException(http.StatusInternalServerError, exception.Enum.InternalServerError))
		return
	}

	err := uc.userService.ForgotPassword(body.Email)
	if err != nil {
		HTTPError(w, r, err.(exception.Exception))
		return
	}

	Respond(w, r, http.StatusOK, nil)
}
