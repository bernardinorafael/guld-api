package org

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s svc) CreateOrg(ctx context.Context, name, ownerId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	org, err := NewOrg(name, ownerId)
	if err != nil {
		msg := "failed to validate org entity"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewValidationFieldError(msg, err, nil)
	}

	if err := s.repo.Insert(ctx, org.Store()); err != nil {
		msg := "failed to create org"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewConflictError(msg, ResourceAlreadyTaken, err, nil)
	}

	return nil
}

func (s svc) GetOrgByID(ctx context.Context, orgId string) (*EntityWithOwner, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	org, err := s.repo.FindByID(ctx, orgId)
	if err != nil {
		msg := "failed to retrieve org"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	return org, nil
}

func (s svc) GetOrgBySlug(ctx context.Context, slug string) (*EntityWithOwner, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	org, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		msg := "failed to retrieve org by slug"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	return org, nil
}
