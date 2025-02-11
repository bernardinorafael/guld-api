package role

import (
	"context"
	"errors"

	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s svc) MakeRolePermissions(ctx context.Context, roleId string, permissions []string) error {
	// TODO: Find and validate team by teamId
	if len(permissions) == 0 {
		msg := "permissions slice cannot be empty"
		s.log.Errorf(ctx, msg, "role_id", roleId)
		return errors.New(msg)
	}

	err := s.repo.BatchRolePermissions(ctx, roleId, permissions)
	if err != nil {
		msg := "failed to insert permissions"
		s.log.Errorf(ctx, msg, "error", err.Error(), "role_id", roleId)
		return errors.New(msg)
	}

	return nil
}

func (s svc) FindByID(ctx context.Context, teamId string, roleId string) (*Entity, error) {
	// TODO: Find and validate team by teamId
	role, err := s.repo.FindByID(ctx, teamId, roleId)
	if err != nil {
		msg := "failed to retrieve role"
		s.log.Errorf(ctx, msg, "error", err.Error(), "role_id", roleId)
		return nil, errors.New(msg)
	}

	return role, nil
}

func (s svc) GetAll(ctx context.Context, teamId string) ([]Entity, error) {
	// TODO: Find and validate team by teamId
	roles, err := s.repo.GetAll(ctx, teamId)
	if err != nil {
		msg := "failed to create role"
		s.log.Errorf(ctx, msg, "error", err.Error())
		return nil, errors.New(msg)
	}

	return roles, nil
}

func (s svc) Create(ctx context.Context, params CreateRoleProps) error {
	// TODO: Find and validate team by teamId

	_, err := s.repo.Create(ctx, Entity{
		Name:        params.Name,
		TeamID:      params.TeamID,
		Key:         params.Key,
		Description: params.Description,
	})
	if err != nil {
		msg := "failed to create role"
		s.log.Errorf(ctx, msg, "error", err.Error())
		return errors.New(msg)
	}

	return nil
}
