package role

import (
	"time"
)

type Permission struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Key  string `json:"key" db:"key"`
}

type Entity struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	TeamID      string       `json:"team_id" db:"team_id"`
	Key         string       `json:"key" db:"key"`
	Description string       `json:"description" db:"description"`
	Permissions []Permission `json:"permissions" db:"permissions"`
	Created     time.Time    `json:"created" db:"created"`
	Updated     time.Time    `json:"updated" db:"updated"`
}

type CreateRoleProps struct {
	Name        string `json:"name"`
	TeamID      string `json:"team_id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

type RolePermissionBatch struct {
	RoleID       string `db:"role_id"`
	PermissionID string `db:"permission_id"`
}
