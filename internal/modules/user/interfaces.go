package user

import (
	"context"

	"github.com/bernardinorafael/internal/_shared/dto"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/pagination"
)

type ServiceInterface interface {
	Create(ctx context.Context, dto UserRegisterDTO) (userId string, err error)
	FindByID(ctx context.Context, userId string) (*CompleteEntity, error)
	Delete(ctx context.Context, userId string) error
	GetAll(ctx context.Context, dto dto.SearchParams) (*pagination.Paginated[EntityWithTeam], error)
	ToggleLock(ctx context.Context, userId string) error
	UpdateProfile(ctx context.Context, userId string, dto UpdateProfileDTO) error
	// Emails methods
	FindAllEmails(ctx context.Context, userId string) ([]email.Entity, error)
	SetPrimaryEmail(ctx context.Context, userId, emailId string) error
	FindEmail(ctx context.Context, email string) (*email.Entity, error)
}

type RepositoryInterface interface {
	Delete(ctx context.Context, userId string) error
	FindByID(ctx context.Context, userId string) (*Entity, error)
	FindCompleteByID(ctx context.Context, userId string) (*CompleteEntity, error)
	GetAll(ctx context.Context, params dto.SearchParams) ([]EntityWithTeam, int, error)
	Create(ctx context.Context, user Entity) error
	Update(ctx context.Context, user Entity) error
}
