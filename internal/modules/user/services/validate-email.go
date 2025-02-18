package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	emails "github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) ValidateEmail(ctx context.Context, emailId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	emailAddr, err := s.emailRepo.FindByID(ctx, emailId)
	if err != nil {
		s.log.Errorw(ctx, "error on get email", logger.Err(err))
		return NewBadRequestError("error on get email", err)
	}
	if emailAddr == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return NewBadRequestError("email not found", err)
	}

	if emailAddr.IsPrimary || emailAddr.IsVerified {
		return NewForbiddenError("invalid status for activating", BadRequest, nil)
	}

	isVerified := true
	err = s.emailRepo.Update(
		ctx,
		emails.EmailUpdateParams{
			ID:         emailAddr.ID,
			IsVerified: &isVerified,
		},
	)
	if err != nil {
		return NewBadRequestError("failed to update email", err)
	}

	return nil
}
