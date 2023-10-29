package apiserver

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	relayStore "relay-backend/internal/store"
)

func Start(config *Config) error {
	db := newDB(config.DatabaseUrl)
	defer db.Close()

	store := relayStore.New(db)
	srv := newServer(store)

	fmt.Printf("Server starting at port %s", config.BindAddress)
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
