package user

import (
	"context"

	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/pagination"
)

type ServiceInterface interface {
	Create(ctx context.Context, user UserRegisterParams) error
	FindByID(ctx context.Context, userId string) (*CompleteEntity, error)
	Delete(ctx context.Context, userId string) error
	GetAll(ctx context.Context, params UserSearchParams) (*pagination.Paginated[Entity], error)
	ToggleLock(ctx context.Context, userId string) error

	// Emails methods
	AddEmail(ctx context.Context, dto CreateEmailParams) error
	FindAllEmails(ctx context.Context, userId string) ([]email.AdditionalEmail, error)
	DeleteEmail(ctx context.Context, userId, emailId string) error
	SetPrimaryEmail(ctx context.Context, userId, emailId string) error
}

type RepositoryInterface interface {
	Delete(ctx context.Context, userId string) error
	FindByID(ctx context.Context, userId string) (*CompleteEntity, error)
	GetAll(ctx context.Context, params UserSearchParams) ([]Entity, int, error)
	Create(ctx context.Context, user Entity) error
	Update(ctx context.Context, input PartialEntity) error
}
