package team

import (
	"context"
	"errors"

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

func (s svc) Create(ctx context.Context, params CreateTeamParams) error {
	// TODO: Check if user exists by `ownerId`
	slug := util.Slugify(params.Name)

	_, err := s.repo.Create(ctx, Entity{
		Name:    params.Name,
		OwnerID: params.OwnerID,
		Slug:    slug,
	})
	if err != nil {
		msg := "failed to create team"
		s.log.Errorf(ctx, msg, "error", err.Error())
		return errors.New(msg)
	}

	return nil
}

func (s svc) GetAll(ctx context.Context, ownerId string) ([]Entity, error) {
	// TODO: Check if user exists by `ownerId`
	teams, err := s.repo.GetAll(ctx, ownerId)
	if err != nil {
		msg := "failed to get teams"
		s.log.Errorf(ctx, msg, "error", err.Error(), "ownerId", ownerId)
		return nil, errors.New(msg)
	}

	return teams, nil
}
