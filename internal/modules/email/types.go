package email

import "time"

type AdditionalEmail struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Email      string    `json:"email" db:"email"`
	IsPrimary  bool      `json:"primary" db:"primary"`
	IsVerified bool      `json:"verified" db:"verified"`
	Created    time.Time `json:"created" db:"created"`
	Updated    time.Time `json:"updated" db:"updated"`
}

type EmailUpdateParams struct {
	ID         string  `json:"id" db:"id"`
	Email      *string `json:"email" db:"email"`
	IsPrimary  *bool   `json:"primary" db:"primary"`
	IsVerified *bool   `json:"verified" db:"verified"`
}
