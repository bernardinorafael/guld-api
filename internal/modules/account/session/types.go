package session

import "time"

type Entity struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	Agent        string    `json:"agent" db:"agent"`
	IP           string    `json:"ip" db:"ip"`
	Revoked      bool      `json:"revoked" db:"revoked"`
	Expires      time.Time `json:"expires" db:"expires"`
	Created      time.Time `json:"created" db:"created"`
	Updated      time.Time `json:"updated" db:"updated"`
}
