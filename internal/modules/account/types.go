package account

import (
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

type Entity struct {
	ID       string    `json:"id" db:"id"`
	UserID   string    `json:"user_id" db:"user_id"`
	Password string    `json:"password,omitempty" db:"password"`
	Created  time.Time `json:"created" db:"created"`
	Updated  time.Time `json:"updated" db:"updated"`
}

type EntityWithUser struct {
	ID       string      `json:"id" db:"id"`
	Password string      `json:"password,omitempty" db:"password"`
	User     user.Entity `json:"user" db:"user"`
	Created  time.Time   `json:"created" db:"created"`
	Updated  time.Time   `json:"updated" db:"updated"`
}

type AccountPayload struct {
	AccessToken string `json:"access_token"`
	AccountID   string `json:"account_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresAt   int64  `json:"expires_at"`
}

type CreateAccountParams struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password,omitempty"`
}
