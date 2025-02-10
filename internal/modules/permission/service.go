package permission

import (
	"context"
	"errors"

	"github.com/bernardinorafael/pkg/logger"
	"github.com/bernardinorafael/pkg/pagination"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s svc) Delete(ctx context.Context, teamId string, permissionId string) error {
	// TODO: validate if team exists
	if err := s.repo.Delete(ctx, teamId, permissionId); err != nil {
		msg := "failed to delete permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "permission_id", permissionId)
		return errors.New(msg)
	}

	return nil
}

func (s svc) Create(ctx context.Context, params CreatePermissionParams) error {
	// TODO: validate if team exists
	err := s.repo.Insert(ctx, Entity{
		Name:        params.Name,
		TeamID:      params.TeamID,
		Key:         params.Key,
		Description: params.Description,
	})
	if err != nil {
		msg := "failed to create permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "team_id", params.TeamID)
		return errors.New(msg)
	}

	return nil
}

func (s svc) GetAll(ctx context.Context, teamId string, p PermissionSearchParams) (*pagination.Paginated[Entity], error) {
	// TODO: validate if team exists
	permissions, count, err := s.repo.GetAll(ctx, teamId, p)
	if err != nil {
		msg := "failed to retrieve permissions"
		s.log.Errorf(ctx, msg, "error", err.Error(), "team_id", teamId)
		return nil, errors.New(msg)
	}
	response := pagination.New(permissions, count, p.Page, p.Limit)

	return &response, nil
}
