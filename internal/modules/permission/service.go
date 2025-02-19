package permission

import (
	"context"
	"time"

	"github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
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

func (s svc) Update(ctx context.Context, dto UpdatePermissionParams) error {
	perm, err := s.repo.GetByID(ctx, dto.OrgID, dto.ID)
	if err != nil {
		msg := "failed to update permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "permission_id", dto.ID)
		return errors.NewBadRequestError(msg, err)
	}
	if perm == nil {
		msg := "permission not found"
		s.log.Errorf(ctx, msg, "permission_id", dto.ID)
		return errors.NewBadRequestError(msg, err)
	}

	err = s.repo.Update(ctx, Entity{
		ID:          perm.ID,
		OrgID:       perm.OrgID,
		Name:        dto.Name,
		Key:         dto.Key,
		Description: dto.Description,
	})
	if err != nil {
		msg := "failed to update permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "permission_id", dto.ID)
		return errors.NewBadRequestError(msg, err)
	}

	return nil
}

func (s svc) Delete(ctx context.Context, teamId string, permissionId string) error {
	// TODO: validate if team exists
	if err := s.repo.Delete(ctx, teamId, permissionId); err != nil {
		msg := "failed to delete permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "permission_id", permissionId)
		return errors.NewBadRequestError(msg, err)
	}

	return nil
}

func (s svc) Create(ctx context.Context, dto CreatePermissionParams) error {
	newPerm := Entity{
		ID:          util.GenID("perm"),
		OrgID:       dto.OrgID,
		Name:        dto.Name,
		Key:         dto.Key,
		Description: dto.Description,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	err := s.repo.Insert(ctx, newPerm)
	if err != nil {
		msg := "failed to create permission"
		s.log.Errorf(ctx, msg, "error", err.Error(), "org_id", dto.OrgID)
		return errors.NewBadRequestError(msg, err)
	}

	return nil
}

func (s svc) GetAll(ctx context.Context, teamId string, dto PermissionSearchParams) (*pagination.Paginated[Entity], error) {
	// TODO: validate if team exists
	permissions, count, err := s.repo.GetAll(ctx, teamId, dto)
	if err != nil {
		msg := "failed to retrieve permissions"
		s.log.Errorf(ctx, msg, "error", err.Error(), "team_id", teamId)
		return nil, errors.NewBadRequestError(msg, err)
	}
	response := pagination.New(permissions, count, dto.Page, dto.Limit)

	return &response, nil
}
