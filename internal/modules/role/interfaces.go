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
	ManagePermissions(ctx context.Context, roleId string, permissions []string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateRoleProps) error
	FindAll(ctx context.Context, orgId string, dto dto.SearchParams) (*pagination.Paginated[EntityWithPermission], error)
	GetRole(ctx context.Context, orgId, roleId string) (*EntityWithPermission, error)
	ManagePermissions(ctx context.Context, roleId string, permissions []string) error
	Delete(ctx context.Context, orgId, roleId string) error
	UpdateRoleInformation(ctx context.Context, orgId, roleId string, dto UpdateRoleDTO) error
}
