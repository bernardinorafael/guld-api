package team

import (
	"context"

	"github.com/bernardinorafael/internal/_shared/dto"
	"github.com/bernardinorafael/pkg/pagination"
)

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
	FindMembersByTeamID(ctx context.Context, orgId, teamId string, dto dto.SearchParams) ([]UserWithRole, int, error)
	DeleteMember(ctx context.Context, userId, teamId string) error
}

type ServiceInterface interface {
	Create(ctx context.Context, dto CreateTeamDTO) error
	GetAll(ctx context.Context, orgId, ownerId string) ([]Entity, error)
	GetBySlug(ctx context.Context, orgId, slug string) (*Entity, error)
	GetByID(ctx context.Context, orgId, teamId string) (*Entity, error)
	AddMember(ctx context.Context, params AddMemberParams) error
	GetByMember(ctx context.Context, orgId, userId string) (*EntityWithRole, error)
	GetMembersByTeamID(ctx context.Context, orgId, teamId string, dto dto.SearchParams) (*pagination.Paginated[UserWithRole], error)
	DeleteMember(ctx context.Context, orgId, userId, teamId string) error
}
