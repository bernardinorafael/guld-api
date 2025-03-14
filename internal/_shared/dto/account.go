package dto

import "time"

type AccountResponse struct {
	SessionID           string `json:"session_id"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenExpires  int64  `json:"access_token_expires"`
	RefreshTokenExpires int64  `json:"refresh_token_expires"`
}

type RenewAccessToken struct {
	AccessToken        string `json:"access_token"`
	AccessTokenExpires int64  `json:"access_token_expires"`
}

type CreateAccount struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password,omitempty"`
}

type SessionResponse struct {
	ID               string    `json:"id"`
	Agent            string    `json:"agent"`
	IP               string    `json:"ip"`
	Revoked          bool      `json:"revoked"`
	Expired          bool      `json:"expired"`
	IsCurrentSession bool      `json:"is_current_session"`
	Created          time.Time `json:"created"`
}
