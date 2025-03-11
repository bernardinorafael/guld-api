package account

import (
	"errors"
	"time"
	"unicode"

	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/user"
)

const (
	minPasswordLength = 6
)

var (
	ErrInvalidPassword = errors.New("failed to encrypt password")
)

type account struct {
	id       string
	userId   string
	password string
	isActive bool
	created  time.Time
	updated  time.Time
}

// NewFromDatabase creates a new account from an entity
func NewFromDatabase(acc Entity) (*account, error) {
	account := &account{
		id:       acc.ID,
		userId:   acc.UserID,
		password: acc.Password,
		created:  acc.Created,
		isActive: acc.IsActive,
		updated:  acc.Updated,
	}

	if err := account.validate(); err != nil {
		return nil, err
	}

	return account, nil
}

// NewAccount creates a new account from scratch
func NewAccount(userId, password string) (*account, error) {
	account := &account{
		id:       util.GenID("acc"),
		userId:   userId,
		password: password,
		isActive: false,
		created:  time.Now(),
		updated:  time.Now(),
	}

	if err := account.validate(); err != nil {
		return nil, err
	}

	return account, nil
}

func (a *account) ChangePassword(password string, ignorePasswordPolicy bool) error {
	// Validate password if ignorePasswordPolicy is false
	if !ignorePasswordPolicy {
		err := a.validatePassword()
		if err != nil {
			return err
		}
	}

	a.password = password
	a.updated = time.Now()

	return nil
}

func (a *account) validate() error {
	if err := a.validatePassword(); err != nil {
		return err
	}

	return nil
}

func (a *account) Activate() {
	a.isActive = true
	a.updated = time.Now()
}

func (a *account) validatePassword() error {
	if a.password == "" {
		return errors.New("password cannot be empty")
	}

	if len(a.password) < minPasswordLength {
		return errors.New("password must be at least 6 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range a.password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func (a *account) StoreWithUser(user user.Entity) EntityWithUser {
	return EntityWithUser{
		ID:       a.ID(),
		Password: a.Password(),
		User:     user,
		Created:  a.Created(),
		Updated:  a.Updated(),
	}
}

func (a *account) Store() Entity {
	return Entity{
		ID:       a.ID(),
		UserID:   a.UserID(),
		Password: a.Password(),
		IsActive: a.IsActive(),
		Created:  a.Created(),
		Updated:  a.Updated(),
	}
}

func (a *account) ID() string         { return a.id }
func (a *account) UserID() string     { return a.userId }
func (a *account) IsActive() bool     { return a.isActive }
func (a *account) Password() string   { return a.password }
func (a *account) Created() time.Time { return a.created }
func (a *account) Updated() time.Time { return a.updated }
