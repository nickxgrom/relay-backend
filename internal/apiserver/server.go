package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"relay-backend/internal/controller"
	"relay-backend/internal/store"
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
	userController := controller.NewUserController(s.store)
	sessionController := controller.NewSessionController(s.store, s.sessionStore)

	s.router.HandleFunc("/users", userController.HandleFunc()).Methods("POST")
	s.router.HandleFunc("/sessions", sessionController.HandleFunc()).Methods("POST")
}
