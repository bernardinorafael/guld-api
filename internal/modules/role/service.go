package role

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bernardinorafael/internal/_shared/dto"
	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/bernardinorafael/pkg/pagination"
	"github.com/lib/pq"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s svc) ManagePermissions(ctx context.Context, roleId string, permissions []string) error {
	err := s.repo.ManagePermissions(ctx, roleId, permissions)
	if err != nil {
		s.log.Errorw(ctx, "failed to manage permissions", logger.Err(err))
		return NewBadRequestError("failed to manage permissions", err)
	}

	return nil
}

func (s svc) GetRole(ctx context.Context, orgId, roleId string) (*EntityWithPermission, error) {
	role, err := s.repo.FindByID(ctx, orgId, roleId)
	if err != nil {
		s.log.Errorw(ctx, "failed to get role", logger.Err(err))
		return nil, NewBadRequestError("failed to get role", err)
	}
	if role == nil {
		return nil, NewNotFoundError("role not found", nil)
	}

	return role, nil
}

func (s *svc) Create(ctx context.Context, dto CreateRoleProps) error {
	role := Entity{
		ID:          util.GenID("role"),
		Name:        dto.Name,
		OrgID:       dto.OrgID,
		Description: dto.Description,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	err := s.repo.Insert(ctx, role)
	if err != nil {
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			var appErr = NewConflictError("role already exists", ResourceAlreadyTaken, err, nil)
			field := util.ExtractFieldFromDetail(pqErr.Detail)
			s.log.Errorw(ctx, "failed to create role already exists", logger.Err(err), logger.Any("field", field))
			appErr.AddField(field, field+" already exists")
			return appErr
		}

		s.log.Errorw(ctx, "failed to create role", logger.Err(err))
		return NewBadRequestError("failed to create role", nil)
	}

	return nil
}

func (s *svc) FindAll(
	ctx context.Context,
	orgId string,
	dto dto.SearchParams,
) (*pagination.Paginated[EntityWithPermission], error) {
	safeSort := map[string]bool{
		"name":    true,
		"created": true,
	}
	// Ignoring `-` preffix on verify sort opts
	sort := strings.TrimPrefix(dto.Sort, "-")
	if !safeSort[sort] {
		s.log.Error(ctx, "invalid sort params")
		return nil, NewValidationFieldError("invalid sort params", nil, []Field{
			{Field: "sort", Msg: "invalid sort params"},
		})
	}

	roles, totalItems, err := s.repo.FindAll(ctx, orgId, dto)
	if err != nil {
		s.log.Errorw(ctx, "failed to find all roles", logger.Err(err))
		return nil, NewBadRequestError("failed to find all roles", nil)
	}

	paginated := pagination.New(roles, totalItems, dto.Page, dto.Limit)
	return &paginated, nil
}

func (s *svc) FindByID(ctx context.Context, orgId string, roleId string) (*Entity, error) {
	panic("unimplemented")
}

func (s *svc) Delete(ctx context.Context, orgId, roleId string) error {
	err := s.repo.Delete(ctx, orgId, roleId)
	if err != nil {
		s.log.Errorw(ctx, "failed to delete role", logger.Err(err))
		return NewBadRequestError("failed to delete role", nil)
	}

	return nil
}

func (s *svc) UpdateRoleInformation(ctx context.Context, orgId, roleId string, dto UpdateRoleDTO) error {
	role, err := s.repo.FindByID(ctx, orgId, roleId)
	if err != nil {
		return NewBadRequestError("failed to update role information", nil)
	}
	if role == nil {
		return NewNotFoundError("role not found", nil)
	}

	err = s.repo.Update(ctx, Entity{
		ID:          role.ID,
		OrgID:       role.OrgID,
		Name:        dto.Name,
		Description: dto.Description,
		Created:     role.Created,
		Updated:     time.Now(),
	})
	if err != nil {
		return NewBadRequestError("failed to update role information", nil)
	}

	return nil
}
