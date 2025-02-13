package team

import "context"

type RepositoryInterface interface {
	Insert(ctx context.Context, team Entity) error
	Delete(ctx context.Context, ownerId, teamId string) error
	FindAll(ctx context.Context, orgId, ownerId string) ([]Entity, error)
	FindAllWithOrg(ctx context.Context, orgId, ownerId string) ([]EntityWithOrg, error)
	FindBySlug(ctx context.Context, orgId, slug string) (*Entity, error)
	FindByID(ctx context.Context, orgId, teamId string) (*Entity, error)
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateTeamParams) error
	GetAll(ctx context.Context, orgId, ownerId string) ([]Entity, error)
	GetBySlug(ctx context.Context, orgId, slug string) (*Entity, error)
	GetByID(ctx context.Context, orgId, teamId string) (*Entity, error)
}
