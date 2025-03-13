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

type AccountPayload struct {
	SessionID           string `json:"session_id"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

type RenewAccessTokenPayload struct {
	AccessToken        string `json:"access_token"`
	AccessTokenExpires int64  `json:"access_token_expires"`
}

type CreateAccountParams struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password,omitempty"`
}
