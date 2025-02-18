package user

import (
	"context"

	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/phone"
	"github.com/bernardinorafael/pkg/pagination"
)

type ServiceInterface interface {
	Create(ctx context.Context, user UserRegisterParams) (userId string, err error)
	FindByID(ctx context.Context, userId string) (*CompleteEntity, error)
	Delete(ctx context.Context, userId string) error
	GetAll(ctx context.Context, params UserSearchParams) (*pagination.Paginated[Entity], error)
	ToggleLock(ctx context.Context, userId string) error

	// Emails methods
	AddEmail(ctx context.Context, dto CreateEmailParams) error
	FindAllEmails(ctx context.Context, userId string) ([]email.AdditionalEmail, error)
	DeleteEmail(ctx context.Context, userId, emailId string) error
	SetPrimaryEmail(ctx context.Context, userId, emailId string) error
	FindEmail(ctx context.Context, email string) (*email.AdditionalEmail, error)
	RequestEmailValidation(ctx context.Context, email, userId string) error
	ValidateEmail(ctx context.Context, emailId string) error

	// Phones methods
	AddPhone(ctx context.Context, dto phone.CreatePhoneParams) error
	FindAllPhones(ctx context.Context, userId string) ([]phone.AdditionalPhone, error)
	DeletePhone(ctx context.Context, userId, phoneId string) error
	SetPrimaryPhone(ctx context.Context, userId, phoneId string) error
}

type RepositoryInterface interface {
	Delete(ctx context.Context, userId string) error
	FindByID(ctx context.Context, userId string) (*Entity, error)
	FindCompleteByID(ctx context.Context, userId string) (*CompleteEntity, error)
	GetAll(ctx context.Context, params UserSearchParams) ([]Entity, int, error)
	Create(ctx context.Context, user Entity) error
	Update(ctx context.Context, input PartialEntity) error
}
