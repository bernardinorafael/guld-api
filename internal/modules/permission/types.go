package permission

import "time"

type Entity struct {
	ID          string    `json:"id" db:"id"`
	OrgID       string    `json:"org_id" db:"org_id"`
	Name        string    `json:"name" db:"name"`
	Key         string    `json:"key" db:"key"`
	Description string    `json:"description" db:"description"`
	Created     time.Time `json:"created" db:"created"`
	Updated     time.Time `json:"updated" db:"updated"`
}

type UpdatePermissionParams struct {
	ID          string `json:"id" db:"id"`
	OrgID       string `json:"org_id" db:"org_id"`
	Name        string `json:"name" db:"name"`
	Key         string `json:"key" db:"key"`
	Description string `json:"description" db:"description"`
}

type CreatePermissionParams struct {
	Name        string `json:"name" db:"name"`
	Key         string `json:"key" db:"key"`
	OrgID       string `json:"org_id" db:"org_id"`
	Description string `json:"description" db:"description"`
}

type PermissionSearchParams struct {
	Sort  string `json:"sort"`
	Query string `json:"q"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
