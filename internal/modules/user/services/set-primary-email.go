package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) SetPrimaryEmail(ctx context.Context, userId, emailId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	primaryEmail, err := s.emailRepo.FindPrimary(ctx, userId)
	if err != nil {
		msg := "error on find primary email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	primaryEmail.IsPrimary = false
	err = s.emailRepo.Update(ctx, email.EmailUpdateParams{
		ID:        primaryEmail.ID,
		IsPrimary: &primaryEmail.IsPrimary,
	})
	if err != nil {
		msg := "error on update primary email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	newPrimaryEmail, err := s.emailRepo.FindByID(ctx, emailId)
	if err != nil {
		msg := "error on find email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	newPrimaryEmail.IsPrimary = true
	err = s.emailRepo.Update(ctx, email.EmailUpdateParams{
		ID:        newPrimaryEmail.ID,
		IsPrimary: &newPrimaryEmail.IsPrimary,
	})
	if err != nil {
		msg := "error on update email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
