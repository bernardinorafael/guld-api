package role

import "context"

type RepositoryInterface interface {
	Insert(ctx context.Context, entity Entity) error
	Update(ctx context.Context, entity Entity) error
	Delete(ctx context.Context, orgId, roleId string) error
	FindByID(ctx context.Context, orgId, roleId string) (*EntityWithPermission, error)
	FindAll(ctx context.Context, orgId string) ([]EntityWithPermission, error)
	// BatchRolePermissions(ctx context.Context, roleId string, permissions []string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateRoleProps) error
	FindByID(ctx context.Context, orgId, roleId string) (*Entity, error)
	FindAll(ctx context.Context, orgId string) ([]EntityWithPermission, error)
}
