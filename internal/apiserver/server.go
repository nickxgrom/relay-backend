package apiserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"relay-backend/internal/controller"
	"relay-backend/internal/store"
)

var (
	sessionName = "auth"
)

type server struct {
	router       *chi.Mux
	store        *store.Store
	sessionStore *sessions.CookieStore
}

func newServer(store *store.Store, sessionStore *sessions.CookieStore) *server {
	s := &server{
		router:       chi.NewRouter(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	authMiddleware := controller.ConfigureMiddleware(s.sessionStore, sessionName)

	s.router.Route("/users", controller.NewUserController(s.store, authMiddleware))
	s.router.Route("/sessions", controller.NewSessionController(s.store, s.sessionStore))
	s.router.Route("/organizations", controller.NewOrganizationController(s.store, authMiddleware))

	//authMiddleware := controller.ConfigureMiddleware(s.sessionStore, sessionName)
	//
	//sessionController := controller.NewSessionController(s.store, s.sessionStore)
	//
	//s.router.HandleFunc("/users", controller.UserHandleFunc(s.store)).Methods(http.MethodPost)
	//s.router.HandleFunc("/sessions", sessionController.HandleFunc()).Methods(http.MethodPost)
	//
	//userRouter := s.router.PathPrefix("").Subrouter()
	//userRouter.Use(authMiddleware.AuthenticateUser)
	//userRouter.HandleFunc("/users", controller.UserHandleFunc(s.store)).Methods(http.MethodGet)
	//
	//organizationRouter := s.router.PathPrefix("/organizations").Subrouter()
	//organizationRouter.Use(authMiddleware.AuthenticateUser)
	//organizationRouter.HandleFunc("/{organizationId}", controller.OrganizationHandleFunc(s.store)).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)
}
