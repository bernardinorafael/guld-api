package permission

import "context"

type RepositoryInterface interface {
	FindAll(ctx context.Context) ([]Entity, error)
	FindByRoleID(ctx context.Context, roleId string) ([]Entity, error)
}

type ServiceInterface interface {
	FindAll(ctx context.Context) ([]Entity, error)
	GetPermissionsByRoleID(ctx context.Context, roleId string) ([]Entity, error)
}
