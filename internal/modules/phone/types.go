package phone

import "time"

type AdditionalPhone struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Phone      string    `json:"phone" db:"phone"`
	IsPrimary  bool      `json:"is_primary" db:"is_primary"`
	IsVerified bool      `json:"is_verified" db:"is_verified"`
	Created    time.Time `json:"created" db:"created"`
	Updated    time.Time `json:"updated" db:"updated"`
}

type CreatePhoneParams struct {
	Phone     string `json:"phone"`
	UserID    string `json:"user_id"`
	IsPrimary bool   `json:"is_primary"`
}

type PhoneUpdateParams struct {
	ID         string `json:"id" db:"id"`
	IsPrimary  *bool  `json:"is_primary" db:"is_primary"`
	IsVerified *bool  `json:"is_verified" db:"is_verified"`
}
