package session

import (
	"context"

	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	log  logger.Logger
	repo RepositoryInterface
}

func NewService(log logger.Logger, repo RepositoryInterface) ServiceInterface {
	return &svc{log, repo}
}

func (s *svc) Create(ctx context.Context, entity Entity) error {
	panic("unimplemented")
}

func (s *svc) Delete(ctx context.Context, sessionId string) error {
	panic("unimplemented")
}

func (s *svc) FindAll(ctx context.Context, username string) ([]Entity, error) {
	panic("unimplemented")
}

func (s *svc) GetSession(ctx context.Context, username string, sessionId string) (*Entity, error) {
	panic("unimplemented")
}
