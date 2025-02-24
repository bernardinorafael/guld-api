package team

import (
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

type Role struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type EntityWithRole struct {
	Entity
	Role Role `json:"role" db:"role"`
}

type UserWithRole struct {
	user.Entity
	Role Role `json:"role"`
}

type Entity struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	OwnerID      string    `json:"owner_id" db:"owner_id"`
	OrgID        string    `json:"org_id" db:"org_id"`
	Logo         *string   `json:"logo" db:"logo"`
	MembersCount int       `json:"members_count" db:"members_count"`
	Created      time.Time `json:"created" db:"created"`
	Updated      time.Time `json:"updated" db:"updated"`
}

type EntityWithMembers struct {
	Entity
	Members []user.Entity `json:"members"`
}
