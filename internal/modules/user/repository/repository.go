package userrepo

import (
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }

func New(db *sqlx.DB) user.RepositoryInterface {
	return &repo{db}
}
