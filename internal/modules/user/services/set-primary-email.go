package usersvc

import (
	"context"
	"fmt"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

/*
* Implementation details
*
* Initially there were methods where we went to the database multiple times
* to fetch data like (previous email and next one to be updated)
* however since there is a limited number of emails that a user can have
* we will do all these operations and validations by fetching all emails at once
 */

func (s svc) SetPrimaryEmail(ctx context.Context, userId, emailId string) error {
	emails, err := s.emailService.FindAllByUser(ctx, userId)
	if err != nil {
		return err
	}

	var currPrimary *email.Entity
	var nextPrimary *email.Entity

	for _, v := range emails {
		if v.IsPrimary {
			currPrimary = &v
		}
		if v.ID == emailId {
			nextPrimary = &v
		}
	}

	if nextPrimary == nil {
		msg := fmt.Sprintf("not found email with id %s", emailId)
		s.log.Errorw(ctx, msg, nil)
		return NewNotFoundError(msg, nil)
	}

	if !nextPrimary.IsVerified {
		return NewForbiddenError("email not verified", EmailNotVerified, nil)
	}

	emailsToUpdate := []email.Entity{
		{ID: currPrimary.ID, IsPrimary: false, IsVerified: currPrimary.IsVerified},
		{ID: nextPrimary.ID, IsPrimary: true, IsVerified: nextPrimary.IsVerified},
	}

	for _, v := range emailsToUpdate {
		_, err = s.emailService.Update(ctx, v)
		if err != nil {
			return err
		}
	}

	err = s.userRepo.Update(ctx, user.PartialEntity{
		ID:           userId,
		EmailAddress: &nextPrimary.Email,
	})
	if err != nil {
		msg := "error on updating user email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
