package team

import "time"

type Entity struct {
	ID      string    `json:"id" db:"id"`
	Name    string    `json:"name" db:"name"`
	Slug    string    `json:"slug" db:"slug"`
	OwnerID string    `json:"owner_id" db:"owner_id"`
	Logo    *string   `json:"logo" db:"logo"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}

type CreateTeamParams struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}
