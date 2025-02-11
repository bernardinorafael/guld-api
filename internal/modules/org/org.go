package org

import (
	"errors"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

type Org struct {
	id        string
	name      string
	ownerID   string
	slug      string
	avatarURL *string
	created   time.Time
	updated   time.Time
}

func NewOrg(name, ownerId string) (*Org, error) {
	org := Org{
		id:        util.GenID("org"),
		slug:      util.Slugify(name),
		name:      name,
		ownerID:   ownerId,
		avatarURL: nil,
		created:   time.Now(),
		updated:   time.Now(),
	}

	if err := org.validate(); err != nil {
		return nil, err
	}

	return &org, nil
}

func (o *Org) validate() error {
	if o.name == "" {
		return errors.New("name is required")
	}

	if o.ownerID == "" {
		return errors.New("ownerID is required")
	}

	return nil
}

func (o *Org) Store() Entity {
	return Entity{
		ID:        o.ID(),
		Name:      o.Name(),
		OwnerID:   o.OwnerID(),
		Slug:      o.Slug(),
		AvatarURL: o.AvatarURL(),
		Created:   o.Created(),
		Updated:   o.Updated(),
	}
}

func (o *Org) ID() string         { return o.id }
func (o *Org) Name() string       { return o.name }
func (o *Org) OwnerID() string    { return o.ownerID }
func (o *Org) Slug() string       { return o.slug }
func (o *Org) AvatarURL() *string { return o.avatarURL }
func (o *Org) Created() time.Time { return o.created }
func (o *Org) Updated() time.Time { return o.updated }
