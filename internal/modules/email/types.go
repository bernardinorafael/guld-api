package email

import "time"

type Entity struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Email      string    `json:"email" db:"email"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Created    time.Time `json:"created" db:"created"`
	Updated    time.Time `json:"updated" db:"updated"`
}

type Validation struct {
	ID         string    `json:"id" db:"id"`
	EmailID    string    `json:"email_id" db:"email_id"`
	Attempts   int       `json:"attempts" db:"attempts"`
	IsConsumed bool      `json:"is_consumed" db:"is_consumed"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Created    time.Time `json:"created" db:"created"`
	Expires    time.Time `json:"expires" db:"expires"`
}
