package email

import (
	"errors"
	"regexp"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

type Email struct {
	id         string
	userId     string
	email      string
	isPrimary  bool
	isVerified bool
	created    time.Time
	updated    time.Time
}

// NewFromEntity creates a new email from an existing entity
func NewFromEntity(e *Entity) (*Email, error) {
	email := &Email{
		id:         e.ID,
		userId:     e.UserID,
		email:      e.Email,
		isPrimary:  e.IsPrimary,
		isVerified: e.IsVerified,
		created:    e.Created,
		updated:    e.Updated,
	}

	if err := email.validate(); err != nil {
		return nil, err
	}

	return email, nil
}

// NewEmail creates a new email entity from scratch
func NewEmail(userId, email string) (*Email, error) {
	e := &Email{
		id:         util.GenID("email"),
		userId:     userId,
		email:      email,
		isPrimary:  false,
		isVerified: false,
		created:    time.Now(),
		updated:    time.Now(),
	}

	if err := e.validate(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Email) Verify() {
	e.isVerified = true
	e.updated = time.Now()
}

func (e *Email) Store() Entity {
	return Entity{
		ID:         e.ID(),
		UserID:     e.UserID(),
		Email:      e.Email(),
		IsPrimary:  e.IsPrimary(),
		IsVerified: e.IsVerified(),
		Created:    e.Created(),
		Updated:    e.Updated(),
	}
}

func (e *Email) ID() string         { return e.id }
func (e *Email) UserID() string     { return e.userId }
func (e *Email) Email() string      { return e.email }
func (e *Email) IsPrimary() bool    { return e.isPrimary }
func (e *Email) IsVerified() bool   { return e.isVerified }
func (e *Email) Created() time.Time { return e.created }
func (e *Email) Updated() time.Time { return e.updated }

func (e *Email) validate() error {
	if e.email == "" {
		return errors.New("email is required")
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(e.email) {
		return errors.New("invalid email format")
	}

	return nil
}
