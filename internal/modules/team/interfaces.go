package team

import "context"

type RepositoryInterface interface {
	Update(ctx context.Context, team Entity) error
	Insert(ctx context.Context, team Entity) error
	Delete(ctx context.Context, ownerId, teamId string) error
	FindAll(ctx context.Context, orgId, ownerId string) ([]Entity, error)
	FindAllWithMembers(ctx context.Context, orgId, ownerId string) ([]EntityWithMembers, error)
	FindBySlug(ctx context.Context, orgId, slug string) (*Entity, error)
	FindByID(ctx context.Context, orgId, teamId string) (*Entity, error)
	InsertMember(ctx context.Context, member TeamMember) error
	FindByMember(ctx context.Context, orgId, userId string) (*EntityWithRole, error)
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateTeamParams) error
	GetAll(ctx context.Context, orgId, ownerId string) ([]Entity, error)
	GetBySlug(ctx context.Context, orgId, slug string) (*Entity, error)
	GetByID(ctx context.Context, orgId, teamId string) (*Entity, error)
	AddMember(ctx context.Context, params AddMemberParams) error
	GetByMember(ctx context.Context, orgId, userId string) (*EntityWithRole, error)
}
