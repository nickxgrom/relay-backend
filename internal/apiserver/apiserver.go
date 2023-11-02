package apiserver

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"relay-backend/internal/store"
)

func Start(config *Config) error {
	db := newDB(config.DatabaseUrl)
	defer db.Close()

	dbStore := store.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(dbStore, sessionStore)

	fmt.Printf("Server starting at port %s\n", config.BindAddress)
	return http.ListenAndServe(config.BindAddress, srv.router)
}

func newDB(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
