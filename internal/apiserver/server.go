package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"relay-backend/internal/controller"
	"relay-backend/internal/model"
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
	userController := controller.NewUserController(s.store)
	sessionController := controller.NewSessionController(s.store, s.sessionStore)

	s.router.HandleFunc("/users", userController.HandleFunc()).Methods("POST")
	s.router.HandleFunc("/sessions", sessionController.HandleFunc()).Methods("POST")

	privateTest := s.router.PathPrefix("/private").Subrouter()
	privateTest.Use(controller.AuthenticateUser(s.sessionStore, userController.GetUserService(), sessionName))
	privateTest.HandleFunc("/whoami", s.whoami).Methods("GET")
}

func (s *server) whoami(w http.ResponseWriter, r *http.Request) {
	controller.Respond(w, r, http.StatusOK, r.Context().Value(controller.CtxKeyUser).(*model.User))
}
