package org

import (
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

type Entity struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	OwnerID   string    `json:"owner_id" db:"owner_id"`
	AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
	Created   time.Time `json:"created" db:"created"`
	Updated   time.Time `json:"updated" db:"updated"`
}

type EntityWithOwner struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	Owner     user.Entity `json:"owner"`
	AvatarURL *string     `json:"avatar_url"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
}
