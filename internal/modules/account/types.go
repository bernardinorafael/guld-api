package account

import (
	"time"

	"github.com/bernardinorafael/internal/modules/org"
	"github.com/bernardinorafael/internal/modules/user"
)

type PartialEntity struct {
	ID       string `db:"id"`
	IsActive *bool  `db:"is_active"`
}

type Entity struct {
	ID       string    `json:"id" db:"id"`
	UserID   string    `json:"user_id" db:"user_id"`
	Password string    `json:"password,omitempty" db:"password"`
	IsActive bool      `json:"is_active" db:"is_active"`
	Created  time.Time `json:"created" db:"created"`
	Updated  time.Time `json:"updated" db:"updated"`
}

type EntityWithUser struct {
	ID       string                  `json:"id" db:"id"`
	Password string                  `json:"password,omitempty" db:"password"`
	IsActive bool                    `json:"is_active" db:"is_active"`
	User     user.Entity             `json:"user" db:"user"`
	Org      *org.EntityWithSettings `json:"org" db:"org"`
	Created  time.Time               `json:"created" db:"created"`
	Updated  time.Time               `json:"updated" db:"updated"`
}

type AccountPayload struct {
	AccessToken string  `json:"access_token"`
	AccountID   string  `json:"account_id"`
	UserID      string  `json:"user_id"`
	OrgID       *string `json:"org_id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	IssuedAt    int64   `json:"issued_at"`
	ExpiresAt   int64   `json:"expires_at"`
}

type CreateAccountParams struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password,omitempty"`
}
