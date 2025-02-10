package permission

import (
	"context"

	"github.com/bernardinorafael/pkg/pagination"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, permission Entity) error
	GetAll(ctx context.Context, teamId string, params PermissionSearchParams) ([]Entity, int, error)
	Delete(ctx context.Context, teamId, permissionId string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreatePermissionParams) error
	GetAll(ctx context.Context, teamId string, params PermissionSearchParams) (*pagination.Paginated[Entity], error)
	Delete(ctx context.Context, teamId, permissionId string) error
}
