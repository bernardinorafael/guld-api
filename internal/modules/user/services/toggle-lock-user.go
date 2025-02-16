package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) ToggleLock(ctx context.Context, userId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	foundUser, err := s.userRepo.FindCompleteByID(ctx, userId)
	if err != nil {
		msg := "failed to retrieve user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewNotFoundError(msg, err)
	}

	usr, err := user.NewFromEntity(foundUser.User)
	if err != nil {
		// TODO: add fields to error
		msg := "error on init user entity"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewValidationFieldError(msg, err, nil)
	}

	if err := usr.ToggleLock(); err != nil {
		msg := "error on toggle lock user status"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	locked := usr.Locked()
	err = s.userRepo.Update(ctx, user.PartialEntity{ID: usr.ID(), Locked: &locked})
	if err != nil {
		msg := "error on toggle lock user status"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil

}
