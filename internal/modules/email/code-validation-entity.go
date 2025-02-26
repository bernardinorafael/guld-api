package email

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bernardinorafael/internal/_shared/util"
)

const (
	MaxAttempts = 3
	CodeLength  = 6
	CodeTTL     = time.Minute * 15
)

type Validation struct {
	id         string
	emailId    string
	userId     string
	code       string
	attempts   int
	isConsumed bool
	isValid    bool
	created    time.Time
	expires    time.Time
}

// NewCodeValidationFromEntity creates a new Validation from a ValidationEntity
func NewCodeValidationFromEntity(entity ValidationEntity) *Validation {
	return &Validation{
		id:         entity.ID,
		emailId:    entity.EmailID,
		userId:     entity.UserID,
		code:       entity.Code,
		attempts:   entity.Attempts,
		isConsumed: entity.IsConsumed,
		isValid:    entity.IsValid,
		created:    entity.Created,
		expires:    entity.Expires,
	}
}

func NewCodeValidation(emailId, userId string) *Validation {
	code := &Validation{
		id:         util.GenID("ev"),
		emailId:    emailId,
		userId:     userId,
		code:       fmt.Sprintf("%06d", rand.Intn(1000000)),
		attempts:   0,
		isConsumed: false,
		isValid:    true,
		created:    time.Now(),
		expires:    time.Now().Add(CodeTTL),
	}

	return code
}

func (v *Validation) ValidateCode(code string) bool {
	if code == v.code {
		v.isValid = true
		v.isConsumed = true
		return true
	}

	return false
}

func (v *Validation) Invalidate() {
	v.isValid = false
}

func (v *Validation) IsMaxAttempts() bool {
	return v.attempts >= MaxAttempts
}

func (v *Validation) IncrementAttempts() {
	v.attempts++
}

func (v *Validation) IsExpired() bool {
	return time.Now().After(v.expires)
}

func (v *Validation) Consume() {
	v.isValid = false
	v.isConsumed = true
}

// Store returns the Entity ready to be stored in the database
func (v *Validation) Store() ValidationEntity {
	return ValidationEntity{
		ID:         v.ID(),
		EmailID:    v.EmailID(),
		UserID:     v.UserID(),
		Code:       v.Code(),
		Attempts:   v.Attempts(),
		IsConsumed: v.IsConsumed(),
		IsValid:    v.IsValid(),
		Created:    v.Created(),
		Expires:    v.Expires(),
	}
}

func (v *Validation) ID() string         { return v.id }
func (v *Validation) EmailID() string    { return v.emailId }
func (v *Validation) UserID() string     { return v.userId }
func (v *Validation) Code() string       { return v.code }
func (v *Validation) Attempts() int      { return v.attempts }
func (v *Validation) IsConsumed() bool   { return v.isConsumed }
func (v *Validation) IsValid() bool      { return v.isValid }
func (v *Validation) Created() time.Time { return v.created }
func (v *Validation) Expires() time.Time { return v.expires }
