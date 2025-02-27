package user

import (
	"time"

	"github.com/bernardinorafael/internal/modules/email"
)

type Team struct {
	ID   *string `json:"id" db:"id"`
	Name *string `json:"name" db:"name"`
}

type EntityWithTeam struct {
	Entity
	Team *Team `json:"team" db:"team,json"`
}

type Entity struct {
	ID           string    `json:"id" db:"id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Username     string    `json:"username" db:"username"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	EmailAddress string    `json:"email_address" db:"email_address"`
	AvatarURL    *string   `json:"avatar_url" db:"avatar_url"`
	Banned       bool      `json:"banned" db:"banned"`
	Locked       bool      `json:"locked" db:"locked"`
	Created      time.Time `json:"created" db:"created"`
	Updated      time.Time `json:"updated" db:"updated"`
}

type CompleteEntity struct {
	User   Entity         `json:"user" db:"user"`
	Emails []email.Entity `json:"emails" db:"emails"`
	Meta   map[string]any `json:"meta"`
}
