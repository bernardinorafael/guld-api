package user

import (
	"context"

	"github.com/bernardinorafael/pkg/pagination"
)

type ServiceInterface interface {
	Create(ctx context.Context, user UserRegisterParams) error
	GetByID(ctx context.Context, userId string) (*CompleteEntity, error)
	Delete(ctx context.Context, userId string) error
	GetAll(ctx context.Context, params UserSearchParams) (*pagination.Paginated[Entity], error)
	ToggleLock(ctx context.Context, userId string) error
}

type RepositoryInterface interface {
	Delete(ctx context.Context, userId string) error
	GetByID(ctx context.Context, userId string) (*CompleteEntity, error)
	GetAll(ctx context.Context, params UserSearchParams) ([]Entity, int, error)
	Create(ctx context.Context, user Entity) error
	Update(ctx context.Context, input PartialEntity) error
}
