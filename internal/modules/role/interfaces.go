package role

import "context"

type RepositoryInterface interface {
	Create(ctx context.Context, role Entity) (*Entity, error)
	Delete(ctx context.Context, teamId, roleId string) error
	GetByID(ctx context.Context, teamId, roleId string) (*Entity, error)
	GetAll(ctx context.Context, teamId string) ([]Entity, error)
	BatchRolePermissions(ctx context.Context, roleId string, permissions []string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateRoleProps) error
	GetByID(ctx context.Context, teamId, roleId string) (*Entity, error)
	GetAll(ctx context.Context, teamId string) ([]Entity, error)
	MakeRolePermissions(ctx context.Context, roleId string, permissions []string) error
}
