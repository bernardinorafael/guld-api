package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) FindByID(ctx context.Context, id string) (*user.CompleteEntity, error) {
	found, err := s.userRepo.FindCompleteByID(ctx, id)
	if err != nil {
		msg := "failed to retrieve user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	return &user.CompleteEntity{
		User: found.User,
		Meta: make([]any, 0),
	}, nil
}
