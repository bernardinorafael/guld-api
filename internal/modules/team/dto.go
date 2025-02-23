package team

import "time"

type CreateTeamParams struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
	OrgID   string `json:"org_id"`
}

type AddMemberParams struct {
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	OrgID  string `json:"org_id"`
}

type TeamMember struct {
	ID      string    `json:"id" db:"id"`
	UserID  string    `json:"user_id" db:"user_id"`
	RoleID  string    `json:"role_id" db:"role_id"`
	TeamID  string    `json:"team_id" db:"team_id"`
	OrgID   string    `json:"org_id" db:"org_id"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}
