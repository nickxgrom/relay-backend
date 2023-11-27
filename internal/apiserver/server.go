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
	authMiddleware := controller.ConfigureMiddleware(s.sessionStore, sessionName, s.store)

	s.router.Route("/users", controller.NewUserController(s.store, authMiddleware))
	s.router.Route("/sessions", controller.NewSessionController(s.store, s.sessionStore))
	s.router.Route("/organizations", controller.NewOrganizationController(s.store, authMiddleware))
}
