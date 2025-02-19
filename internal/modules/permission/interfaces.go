package permission

import (
	"context"

	"github.com/bernardinorafael/pkg/pagination"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, permission Entity) error
	Update(ctx context.Context, permission Entity) error
	GetAll(ctx context.Context, teamId string, params PermissionSearchParams) ([]Entity, int, error)
	GetByID(ctx context.Context, orgId, permId string) (*Entity, error)
	Delete(ctx context.Context, teamId, permId string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, dto CreatePermissionParams) error
	GetAll(ctx context.Context, teamId string, dto PermissionSearchParams) (*pagination.Paginated[Entity], error)
	Delete(ctx context.Context, teamId, permId string) error
	Update(ctx context.Context, dto UpdatePermissionParams) error
}
