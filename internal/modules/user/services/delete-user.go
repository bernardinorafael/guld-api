package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) Delete(ctx context.Context, userId string) error {
	_, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		msg := "user not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewNotFoundError(msg, err)
	}

	if err := s.userRepo.Delete(ctx, userId); err != nil {
		msg := "failed to delete user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
