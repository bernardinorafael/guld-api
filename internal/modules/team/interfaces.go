package team

import "context"

type RepositoryInterface interface {
	Create(ctx context.Context, team Entity) (*Entity, error)
	Delete(ctx context.Context, ownerId, teamId string) error
	GetAll(ctx context.Context, ownerId string) ([]Entity, error)
}

type ServiceInterface interface {
	Create(ctx context.Context, params CreateTeamParams) error
	GetAll(ctx context.Context, ownerId string) ([]Entity, error)
}
