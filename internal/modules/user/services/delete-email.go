package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
)

func (s svc) DeleteEmail(ctx context.Context, userId, emailId string) error {
	emails, err := s.emailService.FindAllByUser(ctx, userId)
	if err != nil {
		return err
	}

	if len(emails) == 1 {
		return NewForbiddenError("cannot delete last user email", InvalidDeletion, err)
	}

	email, err := s.emailService.FindByID(ctx, emailId)
	if err != nil {
		return err
	}

	if email.IsPrimary {
		return NewForbiddenError("primary email cannot be deleted", InvalidDeletion, err)
	}

	if err := s.emailService.Delete(ctx, userId, emailId); err != nil {
		return err
	}

	return nil
}
