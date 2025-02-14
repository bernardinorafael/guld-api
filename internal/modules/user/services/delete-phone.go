package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) DeletePhone(ctx context.Context, userId, phoneId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	allPhones, err := s.phoneRepo.FindAllByUser(ctx, userId)
	if err != nil {
		msg := "error on find all phones"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	if len(allPhones) == 1 {
		msg := "cannot delete last phone"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewForbiddenError(msg, InvalidDeletion, err)
	}

	phone, err := s.phoneRepo.FindByID(ctx, phoneId)
	if err != nil {
		msg := "error on find phone"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	if phone == nil {
		msg := "phone not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	if phone.IsPrimary {
		msg := "primary phone cannot be deleted"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewForbiddenError(msg, InvalidDeletion, err)
	}

	if err := s.phoneRepo.Delete(ctx, userId, phoneId); err != nil {
		msg := "error on delete phone"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
