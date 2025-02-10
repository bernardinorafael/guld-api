package pg

import (
	"context"
	"fmt"

	"github.com/bernardinorafael/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(log logger.Logger, dsn string) (*Database, error) {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		log.Errorf(context.TODO(), "error opening database", "error", err.Error())
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return &Database{db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
