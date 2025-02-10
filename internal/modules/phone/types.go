package phone

import "time"

type AdditionalPhone struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Phone      string    `json:"phone" db:"phone"`
	IsPrimary  bool      `json:"primary" db:"primary"`
	IsVerified bool      `json:"verified" db:"verified"`
	Created    time.Time `json:"created" db:"created"`
	Updated    time.Time `json:"updated" db:"updated"`
}

type PhoneUpdateParams struct {
	ID         string  `json:"id" db:"id"`
	Phone      *string `json:"phone" db:"phone"`
	IsPrimary  *bool   `json:"primary" db:"primary"`
	IsVerified *bool   `json:"verified" db:"verified"`
}
