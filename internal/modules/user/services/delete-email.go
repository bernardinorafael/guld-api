package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) DeleteEmail(ctx context.Context, userId, emailId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	email, err := s.emailRepo.FindByID(ctx, emailId)
	if err != nil {
		msg := "error on find email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	if email == nil {
		msg := "email not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	if email.IsPrimary {
		msg := "primary email cannot be deleted"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewForbiddenError(msg, InvalidDeletion, err)
	}

	if err := s.emailRepo.Delete(ctx, userId, emailId); err != nil {
		msg := "error on delete email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
