package team

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bernardinorafael/internal/_shared/dto"
	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/bernardinorafael/pkg/pagination"
	"github.com/lib/pq"
)

type WithParams struct {
	Log      logger.Logger
	TeamRepo RepositoryInterface
}

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s *svc) DeleteMember(ctx context.Context, orgId, userId, teamId string) error {
	t, err := s.repo.FindByID(ctx, orgId, teamId)
	if err != nil {
		return NewBadRequestError("failed to get team by teamId", err)
	}

	err = s.repo.DeleteMember(ctx, userId, teamId)
	if err != nil {
		return NewBadRequestError("failed to delete member", err)
	}

	team, err := NewFromEntity(*t)
	if err != nil {
		return NewBadRequestError("failed to create team", err)
	}
	team.DecrementMembersCount()

	err = s.repo.Update(ctx, team.Store())
	if err != nil {
		return NewBadRequestError("failed to update team members count", err)
	}

	return nil
}

func (s *svc) GetMembersByTeamID(
	ctx context.Context,
	orgId, slug string,
	input dto.SearchParams,
) (*pagination.Paginated[UserWithRole], error) {
	safeSort := map[string]bool{
		"full_name": true,
		"username":  true,
		"created":   true,
	}
	// Ignoring `-` preffix on verify sort opts
	sort := strings.TrimPrefix(input.Sort, "-")
	if !safeSort[sort] {
		msg := "invalid sort params"
		s.log.Errorw(ctx, msg, logger.Err(fmt.Errorf("invalid sort params: %s", input.Sort)))
		return nil, NewValidationFieldError(
			msg,
			fmt.Errorf("invalid sort params: %s", input.Sort),
			[]Field{
				{Field: "sort", Msg: "invalid sort params"},
			},
		)
	}

	team, err := s.repo.FindBySlug(ctx, orgId, slug)
	if err != nil {
		s.log.Errorw(ctx, "failed to get team by id", logger.Err(err))
		return nil, NewBadRequestError("failed to get team by id", err)
	}
	if team == nil {
		s.log.Errorw(ctx, "team not found", logger.Err(err))
		return nil, NewNotFoundError("team not found", err)
	}

	members, totalItems, err := s.repo.FindMembersByTeamID(ctx, orgId, team.ID, input)
	if err != nil {
		s.log.Errorw(ctx, "failed to get members by team id", logger.Err(err))
		return nil, NewBadRequestError("failed to get members by team id", err)
	}

	paginated := pagination.New(members, totalItems, input.Page, input.Limit)
	return &paginated, nil
}

func (s *svc) GetByMember(ctx context.Context, orgId, userId string) (*EntityWithRole, error) {
	team, err := s.repo.FindByMember(ctx, orgId, userId)
	if err != nil {
		s.log.Errorw(ctx, "failed to get team by member", logger.Err(err))
		return nil, NewBadRequestError("failed to get team by member", err)
	}
	if team == nil {
		s.log.Warn(ctx, "team not found")
	}

	return team, nil
}

func (s *svc) AddMember(ctx context.Context, input AddMemberParams) error {
	t, err := s.repo.FindByID(ctx, input.OrgID, input.TeamID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get team by id", logger.Err(err))
		return NewBadRequestError("failed to get team by id", err)
	}

	team, err := NewFromEntity(*t)
	if err != nil {
		s.log.Errorw(ctx, "failed to create team entity", logger.Err(err))
		return NewBadRequestError("failed to create team", err)
	}

	member := TeamMember{
		ID:      util.GenID("tm"),
		UserID:  input.UserID,
		RoleID:  input.RoleID,
		TeamID:  input.TeamID,
		OrgID:   input.OrgID,
		Created: time.Now(),
		Updated: time.Now(),
	}

	err = s.repo.InsertMember(ctx, member)
	if err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return NewForbiddenError("member already in team", DuplicatedField, err)
		}
		s.log.Errorw(ctx, "failed to add member to team", logger.Err(err))
		return NewBadRequestError("failed to add member to team", err)
	}

	team.IncrementMembersCount()
	err = s.repo.Update(ctx, team.Store())
	if err != nil {
		s.log.Errorw(ctx, "failed to update team members count", logger.Err(err))
		return NewBadRequestError("failed to update team members count", err)
	}

	return nil
}

func (s *svc) GetByID(ctx context.Context, orgId string, teamId string) (*Entity, error) {
	team, err := s.repo.FindByID(ctx, orgId, teamId)
	if err != nil {
		s.log.Errorw(ctx, "failed to get team by id", logger.Err(err))
		return nil, NewBadRequestError("failed to get team by id", err)
	}
	if team == nil {
		s.log.Errorw(ctx, "team not found", logger.Err(err))
		return nil, NewNotFoundError("team not found", err)
	}

	return team, nil
}

func (s svc) GetBySlug(ctx context.Context, orgId string, slug string) (*Entity, error) {
	team, err := s.repo.FindBySlug(ctx, orgId, slug)
	if err != nil {
		s.log.Errorw(ctx, "failed to get team by slug", logger.Err(err))
		return nil, NewBadRequestError("failed to get team by slug", err)
	}
	if team == nil {
		s.log.Errorw(ctx, "team not found", logger.Err(err))
		return nil, NewNotFoundError("team not found", err)
	}

	return team, nil
}

func (s svc) Create(ctx context.Context, dto CreateTeamDTO) error {
	newTeam, err := NewTeam(dto.Name, dto.OwnerID, dto.OrgID)
	if err != nil {
		s.log.Errorw(ctx, "failed to create team", logger.Err(err))
		return NewBadRequestError("failed to create team", err)
	}
	toStore := newTeam.Store()

	err = s.repo.Insert(ctx, toStore)
	if err != nil {
		s.log.Errorw(ctx, "failed to create team", logger.Err(err))
		return NewBadRequestError("failed to create team", err)
	}

	return nil
}

func (s svc) GetAll(ctx context.Context, ownerId, orgId string) ([]Entity, error) {
	teams, err := s.repo.FindAll(ctx, ownerId, orgId)
	if err != nil {
		s.log.Errorw(ctx, "failed to get teams", logger.Err(err))
		return nil, NewBadRequestError("failed to get teams", err)
	}

	return teams, nil
}
