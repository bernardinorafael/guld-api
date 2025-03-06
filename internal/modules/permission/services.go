package permission

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

func (s *svc) FindAll(ctx context.Context) ([]Entity, error) {
	permissions, err := s.repo.FindAll(ctx)
	if err != nil {
		s.log.Errorw(ctx, "failed to find permissions", logger.Err(err))
		return nil, NewBadRequestError("failed to find permissions", err)
	}

	return permissions, nil
}
