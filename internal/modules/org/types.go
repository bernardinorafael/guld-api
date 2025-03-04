package org

import (
	"time"

	"github.com/bernardinorafael/internal/modules/user"
)

type Settings struct {
	ID                        string    `json:"id" db:"id"`
	OrgID                     string    `json:"org_id" db:"org_id"`
	IsActive                  bool      `json:"is_active" db:"is_active"`
	DefaultMembershipPassword string    `json:"default_membership_password" db:"default_membership_password"`
	MaxAllowedMemberships     int       `json:"max_allowed_memberships" db:"max_allowed_memberships"`
	MaxAllowedRoles           int       `json:"max_allowed_roles" db:"max_allowed_roles"`
	UseMasterPassword         bool      `json:"use_master_password" db:"use_master_password"`
	Created                   time.Time `json:"created" db:"created"`
	Updated                   time.Time `json:"updated" db:"updated"`
}

type EntityWithSettings struct {
	Entity
	Settings Settings `json:"settings"`
}

type Entity struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	OwnerID   string    `json:"owner_id" db:"owner_id"`
	AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
	Created   time.Time `json:"created" db:"created"`
	Updated   time.Time `json:"updated" db:"updated"`
}

type EntityWithOwner struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	Owner     user.Entity `json:"owner"`
	AvatarURL *string     `json:"avatar_url"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
}
