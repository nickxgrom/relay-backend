package apiserver

import (
	"github.com/gorilla/mux"
	"relay-backend/internal/controller"
	"relay-backend/internal/store"
)

type server struct {
	router *mux.Router
	store  *store.Store
}

func newServer(store *store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	userController := controller.NewUserController(s.store)
	sessionController := controller.NewSessionController(s.store)

	s.router.HandleFunc("/users", userController.HandleFunc()).Methods("POST")
	s.router.HandleFunc("/sessions", sessionController.HandleFunc()).Methods("POST")
}
