package usersvc

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/email"
	emails "github.com/bernardinorafael/internal/modules/email"
)

func (s svc) ValidateEmail(ctx context.Context, emailId string) error {
	emailAddr, err := s.emailService.FindByID(ctx, emailId)
	if err != nil {
		return err
	}

	if emailAddr.IsPrimary || emailAddr.IsVerified {
		s.log.Errorw(ctx, "invalid status for activating email")
		return NewForbiddenError("invalid status for activating", BadRequest, nil)
	}

	isVerified := true
	update := emails.Entity{ID: emailAddr.ID, IsVerified: isVerified}
	_, err = s.emailService.Update(ctx, update)
	if err != nil {
		return err
	}

	validation, err := s.emailService.FindValidationByEmail(ctx, emailAddr.ID)
	if err != nil {
		return err
	}

	if validation.Expires.Before(time.Now()) {
		s.log.Error(ctx, "email validation expired")
		return NewForbiddenError("email validation expired", ExpiredLink, nil)
	}

	toUpdate := email.Validation{ID: validation.ID, IsConsumed: true}
	err = s.emailService.UpdateValidation(ctx, toUpdate)
	if err != nil {
		return NewBadRequestError("failed to update email validation", err)
	}

	return nil
}
