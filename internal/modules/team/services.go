package team

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
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

func (s *svc) GetByID(ctx context.Context, orgId string, teamId string) (*Entity, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	team, err := s.repo.FindByID(ctx, orgId, teamId)
	if err != nil {
		msg := "failed to get team by id"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}
	if team == nil {
		msg := "team not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	return team, nil
}

func (s svc) GetBySlug(ctx context.Context, orgId string, slug string) (*Entity, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	team, err := s.repo.FindBySlug(ctx, orgId, slug)
	if err != nil {
		msg := "failed to get team by slug"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}
	if team == nil {
		msg := "team not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	return team, nil
}

func (s svc) Create(ctx context.Context, dto CreateTeamParams) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	newTeam, err := NewTeam(dto.Name, dto.OwnerID, dto.OrgID)
	if err != nil {
		msg := "failed to create team"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	toStore := newTeam.Store()

	err = s.repo.Insert(ctx, toStore)
	if err != nil {
		msg := "failed to create team"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}

func (s svc) GetAll(ctx context.Context, ownerId, orgId string) ([]Entity, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	teams, err := s.repo.FindAll(ctx, ownerId, orgId)
	if err != nil {
		msg := "failed to get teams"
		s.log.Errorf(ctx, msg, "error", err.Error(), "ownerId", ownerId)
		return nil, NewBadRequestError(msg, err)
	}

	return teams, nil
}
