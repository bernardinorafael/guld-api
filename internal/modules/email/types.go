package email

import "time"

type AdditionalEmail struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Email      string    `json:"email" db:"email"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Created    time.Time `json:"created" db:"created"`
	Updated    time.Time `json:"updated" db:"updated"`
}

type EmailUpdateParams struct {
	ID         string `json:"id" db:"id"`
	IsPrimary  *bool  `json:"is_primary" db:"is_primary"`
	IsVerified *bool  `json:"is_verified" db:"is_verified"`
}
