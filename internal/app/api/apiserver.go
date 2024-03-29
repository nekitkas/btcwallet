package api

import (
	"btcwallet/internal/store/sqlstore"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseDriver, config.DatabaseURL, config.DatabaseMigrations)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(store)

	srv.logger.Printf("Server started at port %v", config.Port)

	return http.ListenAndServe(config.Port, srv)
}

func newDB(driver, url, migrations string) (*sql.DB, error) {
	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("instance: %w", err)
	}

	fileSource, err := (&file.File{}).Open(migrations)
	if err != nil {
		return nil, fmt.Errorf("fileSource: %w", err)
	}

	m, err := migrate.NewWithInstance("file", fileSource, driver, instance)
	if err != nil {
		return nil, fmt.Errorf("migrations new: %w", err)
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, fmt.Errorf("migrations run: %w", err)
		}
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
