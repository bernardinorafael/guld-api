package team

import (
	"time"

	"github.com/bernardinorafael/internal/modules/org"
)

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

type EntityWithOrg struct {
	Entity
	Org *org.Entity `json:"org"`
}

type CreateTeamParams struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
	OrgID   string `json:"org_id"`
}
