package role

import (
	"context"
	"time"

	"github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s *svc) Create(ctx context.Context, dto CreateRoleProps) error {
	newRole := Entity{
		ID:          util.GenID("role"),
		Name:        dto.Name,
		OrgID:       dto.OrgID,
		Key:         dto.Key,
		Description: dto.Description,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	err := s.repo.Insert(ctx, newRole)
	if err != nil {
		s.log.Errorw(ctx, "failed to create role", logger.Err(err))
		return errors.NewBadRequestError("failed to create role", nil)
	}

	return nil
}

func (s *svc) FindAll(ctx context.Context, orgId string) ([]EntityWithPermission, error) {
	roles, err := s.repo.FindAll(ctx, orgId)
	if err != nil {
		s.log.Errorw(ctx, "failed to find all roles", logger.Err(err))
		return nil, errors.NewBadRequestError("failed to find all roles", nil)
	}

	return roles, nil
}

func (s *svc) FindByID(ctx context.Context, orgId string, roleId string) (*Entity, error) {
	panic("unimplemented")
}
