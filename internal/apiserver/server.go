package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/controller"
	"relay-backend/internal/store"
)

var (
	sessionName = "auth"
)

type server struct {
	router       *mux.Router
	store        *store.Store
	sessionStore *sessions.CookieStore
}

func newServer(store *store.Store, sessionStore *sessions.CookieStore) *server {
	s := &server{
		router:       mux.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	authMiddleware := controller.ConfigureMiddleware(s.sessionStore, sessionName)

	sessionController := controller.NewSessionController(s.store, s.sessionStore)

	s.router.HandleFunc("/users", controller.UserHandleFunc(s.store)).Methods(http.MethodPost)
	s.router.HandleFunc("/sessions", sessionController.HandleFunc()).Methods(http.MethodPost)

	userRouter := s.router.PathPrefix("").Subrouter()
	userRouter.Use(authMiddleware.AuthenticateUser)
	userRouter.HandleFunc("/users", controller.UserHandleFunc(s.store)).Methods(http.MethodGet)

	organizationRouter := s.router.PathPrefix("/organizations").Subrouter()
	organizationRouter.Use(authMiddleware.AuthenticateUser)
	organizationRouter.HandleFunc("", controller.OrganizationHandleFunc(s.store)).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
	organizationRouter.HandleFunc("/list", controller.OrganizationListHandleFunc(s.store)).Methods(http.MethodGet)
}
