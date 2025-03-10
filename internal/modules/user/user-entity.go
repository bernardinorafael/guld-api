package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

const (
	minNameLength = 3
	// username lockout duration - 90 days
	usernameLockoutDuration = 90 * 24 * time.Hour
)

var (
	// Name validation errors
	ErrInvalidFullNameLength = fmt.Errorf("full name must be at least %d characters long", minNameLength)
	ErrInvalidFullName       = errors.New("incorrect name, must contain valid first and last name")
	ErrEmptyFullName         = errors.New("full name is a required field")
	// Username validation errors
	ErrEmptyUsername  = errors.New("username is a required field")
	ErrUsernameLocked = errors.New("username is locked for update")
	// ID validation errors
	ErrInvalidID = errors.New("invalid ksuid format")
)

type User struct {
	id           string
	fullName     string
	username     string
	phoneNumber  string
	emailAddress string
	avatarURL    *string

	banned               bool
	locked               bool
	ignorePasswordPolicy bool

	usernameLastUpdated time.Time
	usernameLockoutEnd  time.Time
	created             time.Time
	updated             time.Time
}

// NewFromEntity creates a new user entity from an existing one
func NewFromEntity(u Entity) (*User, error) {
	user := User{
		id:           u.ID,
		fullName:     u.FullName,
		username:     u.Username,
		phoneNumber:  u.PhoneNumber,
		emailAddress: u.EmailAddress,
		avatarURL:    u.AvatarURL,

		banned:               u.Banned,
		locked:               u.Locked,
		ignorePasswordPolicy: u.IgnorePasswordPolicy,

		usernameLastUpdated: u.UsernameLastUpdated,
		usernameLockoutEnd:  u.UsernameLockoutEnd,
		created:             u.Created,
		updated:             u.Updated,
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return &user, nil
}

// NewUser creates a new user entity from scratch
func NewUser(name, username, phone, email string) (*User, error) {

	user := User{
		id:           util.GenID("user"),
		fullName:     name,
		username:     username,
		phoneNumber:  phone,
		emailAddress: email,
		avatarURL:    nil,

		banned:               false,
		locked:               false,
		ignorePasswordPolicy: false,

		usernameLastUpdated: time.Now(),
		usernameLockoutEnd:  time.Now().Add(usernameLockoutDuration),
		created:             time.Now(),
		updated:             time.Now(),
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) validate() error {
	if u.fullName == "" {
		return ErrEmptyFullName
	}

	if len(u.fullName) < minNameLength {
		return ErrInvalidFullNameLength
	}

	if len(u.username) == 0 {
		return ErrEmptyUsername
	}

	return nil
}

func (u *User) ToggleLock() error {
	if u.locked {
		return u.Unlock()
	}
	return u.Lock()
}

func (u *User) Lock() error {
	if u.locked {
		return fmt.Errorf("user is already locked")
	}

	u.locked = true
	u.updated = time.Now()

	return nil
}

func (u *User) Unlock() error {
	if !u.locked {
		return fmt.Errorf("user is already unlocked")
	}

	u.locked = false
	u.updated = time.Now()

	return nil
}

func (u *User) ChangeEmail(email string) error {
	if err := u.validate(); err != nil {
		return err
	}

	u.emailAddress = email
	u.updated = time.Now()

	return nil
}

func (u *User) ChangePhone(phone string) error {
	if err := u.validate(); err != nil {
		return err
	}

	u.phoneNumber = phone
	u.updated = time.Now()
	return nil
}

func (u *User) ChangeName(name string) error {
	if err := u.validate(); err != nil {
		return err
	}

	u.fullName = name
	u.updated = time.Now()
	return nil
}

func (u *User) ChangeUsername(username string) error {
	// if the username is unchanged, do nothing and skip validation
	if u.username == username {
		return nil
	}

	if u.usernameLockoutEnd.After(time.Now()) {
		return ErrUsernameLocked
	}

	err := u.validate()
	if err != nil {
		return err
	}

	u.username = username
	u.usernameLastUpdated = time.Now()
	u.usernameLockoutEnd = time.Now().Add(usernameLockoutDuration)
	u.updated = time.Now()
	return nil
}

func (u *User) ChangeProfilePicture(url string) {
	u.avatarURL = &url
	u.updated = time.Now()
}

func (u *User) ChangePasswordPolicy(ignore bool) {
	u.ignorePasswordPolicy = ignore
	u.updated = time.Now()
}

// Store stores the user entity in the database
func (u *User) Store() Entity {
	return Entity{
		ID:           u.ID(),
		FullName:     u.FullName(),
		Username:     u.Username(),
		AvatarURL:    u.AvatarURL(),
		PhoneNumber:  u.Phone(),
		EmailAddress: u.Email(),

		Banned:               u.Banned(),
		Locked:               u.Locked(),
		IgnorePasswordPolicy: u.IgnorePasswordPolicy(),

		UsernameLastUpdated: u.UsernameLastUpdated(),
		UsernameLockoutEnd:  u.UsernameLockoutEnd(),
		Created:             u.Created(),
		Updated:             u.Updated(),
	}
}

func (u *User) ID() string                     { return u.id }
func (u *User) FullName() string               { return u.fullName }
func (u *User) Username() string               { return u.username }
func (u *User) Email() string                  { return u.emailAddress }
func (u *User) AvatarURL() *string             { return u.avatarURL }
func (u *User) UsernameLastUpdated() time.Time { return u.usernameLastUpdated }
func (u *User) UsernameLockoutEnd() time.Time  { return u.usernameLockoutEnd }
func (u *User) Phone() string                  { return u.phoneNumber }
func (u *User) Banned() bool                   { return u.banned }
func (u *User) Locked() bool                   { return u.locked }
func (u *User) IgnorePasswordPolicy() bool     { return u.ignorePasswordPolicy }
func (u *User) Created() time.Time             { return u.created }
func (u *User) Updated() time.Time             { return u.updated }
