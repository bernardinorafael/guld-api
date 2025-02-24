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
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	OrgID       string    `json:"org_id" db:"org_id"`
	Description string    `json:"description" db:"description"`
	Created     time.Time `json:"created" db:"created"`
	Updated     time.Time `json:"updated" db:"updated"`
}

type EntityWithPermission struct {
	ID          string       `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	OrgID       string       `json:"org_id" db:"org_id"`
	Permissions []Permission `json:"permissions" db:"permissions"`
	Description string       `json:"description" db:"description"`
	Created     time.Time    `json:"created" db:"created"`
	Updated     time.Time    `json:"updated" db:"updated"`
}

type CreateRoleProps struct {
	Name        string `json:"name"`
	OrgID       string `json:"org_id"`
	Description string `json:"description"`
}

type RolePermissionBatch struct {
	RoleID       string `db:"role_id"`
	PermissionID string `db:"permission_id"`
}
