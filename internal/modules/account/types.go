package account

import (
	"time"

	"github.com/bernardinorafael/internal/modules/org"
	"github.com/bernardinorafael/internal/modules/user"
)

type Entity struct {
	ID       string    `json:"id" db:"id"`
	UserID   string    `json:"user_id" db:"user_id"`
	Password string    `json:"password,omitempty" db:"password"`
	IsActive bool      `json:"is_active" db:"is_active"`
	Created  time.Time `json:"created" db:"created"`
	Updated  time.Time `json:"updated" db:"updated"`
}

// TODO: use embbeded struct
type EntityWithUser struct {
	ID       string                  `json:"id" db:"id"`
	Password string                  `json:"password,omitempty" db:"password"`
	IsActive bool                    `json:"is_active" db:"is_active"`
	User     user.Entity             `json:"user" db:"user"`
	Org      *org.EntityWithSettings `json:"org" db:"org"`
	Created  time.Time               `json:"created" db:"created"`
	Updated  time.Time               `json:"updated" db:"updated"`
}
