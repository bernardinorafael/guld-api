package team

import (
	"errors"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

const (
	minNameLength = 3
)

type Team struct {
	id           string
	name         string
	slug         string
	ownerId      string
	orgId        string
	membersCount int
	logo         *string
	created      time.Time
	updated      time.Time
}

func NewFromEntity(entity Entity) (*Team, error) {
	team := &Team{
		id:           entity.ID,
		name:         entity.Name,
		slug:         entity.Slug,
		ownerId:      entity.OwnerID,
		orgId:        entity.OrgID,
		membersCount: entity.MembersCount,
		logo:         entity.Logo,
		created:      entity.Created,
		updated:      entity.Updated,
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func NewTeam(name, ownerId, orgId string) (*Team, error) {
	team := &Team{
		id:           util.GenID("team"),
		slug:         util.Slugify(name),
		name:         name,
		ownerId:      ownerId,
		orgId:        orgId,
		logo:         nil,
		membersCount: 0,
		created:      time.Now(),
		updated:      time.Now(),
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) IncrementMembersCount() {
	t.membersCount++
	t.updated = time.Now()
}

func (t *Team) DecrementMembersCount() {
	if t.membersCount > 0 {
		t.membersCount--
		t.updated = time.Now()
	}
}

func (t *Team) validate() error {
	if t.id == "" {
		return errors.New("id is a required field")
	}

	if t.name == "" {
		return errors.New("name is a required field")
	}

	if t.ownerId == "" {
		return errors.New("ownerId is a required field")
	}

	if t.orgId == "" {
		return errors.New("orgId is a required field")
	}

	if len(t.name) < minNameLength {
		return errors.New("name must be at least 3 characters long")
	}

	return nil
}

func (t *Team) Store() Entity {
	return Entity{
		ID:           t.ID(),
		Name:         t.Name(),
		Slug:         t.Slug(),
		OwnerID:      t.OwnerID(),
		OrgID:        t.OrgID(),
		MembersCount: t.MembersCount(),
		Logo:         t.Logo(),
		Created:      t.Created(),
		Updated:      t.Updated(),
	}
}

func (t *Team) ID() string         { return t.id }
func (t *Team) Name() string       { return t.name }
func (t *Team) Slug() string       { return t.slug }
func (t *Team) OwnerID() string    { return t.ownerId }
func (t *Team) OrgID() string      { return t.orgId }
func (t *Team) MembersCount() int  { return t.membersCount }
func (t *Team) Logo() *string      { return t.logo }
func (t *Team) Created() time.Time { return t.created }
func (t *Team) Updated() time.Time { return t.updated }
