package role

import (
	"context"

	"github.com/bernardinorafael/internal/_shared/dto"
	"github.com/bernardinorafael/pkg/pagination"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, entity Entity) error
	Update(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, orgId, roleId string) error
	FindByID(ctx context.Context, orgId, roleId string) (*EntityWithPermission, error)
	FindAll(ctx context.Context, orgId string, dto dto.SearchParams) ([]EntityWithPermission, int, error)
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateRoleProps) error
	FindByID(ctx context.Context, orgId, roleId string) (*Entity, error)
	FindAll(ctx context.Context, orgId string, dto dto.SearchParams) (*pagination.Paginated[EntityWithPermission], error)
}
