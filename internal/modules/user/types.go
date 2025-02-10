package user

import (
	"time"

	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/phone"
)

type PartialEntity struct {
	ID string `db:"id"`

	FullName     *string `db:"full_name"`
	Username     *string `db:"username"`
	PhoneNumber  *string `db:"phone_number"`
	EmailAddress *string `db:"email_address"`
	AvatarURL    *string `db:"avatar_url"`
	Banned       *bool   `db:"banned"`
	Locked       *bool   `db:"locked"`
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
	User   Entity                  `json:"user" db:"user"`
	Emails []email.AdditionalEmail `json:"emails" db:"emails"`
	Phones []phone.AdditionalPhone `json:"phones" db:"phones"`
	Meta   []any                   `json:"meta"`
}
