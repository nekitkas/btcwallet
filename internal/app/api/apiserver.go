package api

import (
	"btcwallet/internal/store/sqlstore"
	"database/sql"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseDriver, config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(store)

	srv.logger.Printf("Server started at port %v", config.Port)

	return http.ListenAndServe(config.Port, srv)
}

func newDB(driver, url string) (*sql.DB, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
