package team

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/pkg/logger"
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

func (s *svc) GetTeamsByMember(ctx context.Context, orgId, userId string) ([]Entity, error) {
	teams, err := s.repo.FindTeamsByMember(ctx, orgId, userId)
	if err != nil {
		s.log.Errorw(ctx, "failed to get teams by member", logger.Err(err))
		return nil, NewBadRequestError("failed to get teams by member", err)
	}

	return teams, nil
}

func (s *svc) AddMember(ctx context.Context, input AddMemberParams) error {
	tm, err := s.repo.FindTeamsByMember(ctx, input.OrgID, input.UserID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get teams by member", logger.Err(err))
		return NewBadRequestError("failed to get teams by member", err)
	}

	if len(tm) >= 3 {
		s.log.Errorw(ctx, "max limit of teams reached", logger.Err(err))
		return NewForbiddenError("max limit of teams reached", MaxLimitResourceReached, err)
	}

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

func (s svc) Create(ctx context.Context, dto CreateTeamParams) error {
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
