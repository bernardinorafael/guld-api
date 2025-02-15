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

* Initially there were methods where we went to the database multiple times
* to fetch data like (previous email and next one to be updated)
* however since there is a limited number of emails that a user can have
* we will do all these operations and validations by fetching all emails at once
 */

func (s svc) SetPrimaryEmail(ctx context.Context, userId, emailId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	emails, err := s.emailRepo.FindAllByUser(ctx, userId)
	if err != nil {
		msg := "error on find all emails"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	var currPrimary *email.AdditionalEmail
	var nextPrimary *email.AdditionalEmail

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

	currPrimary.IsPrimary = false
	nextPrimary.IsPrimary = true

	emailsToUpdate := []email.EmailUpdateParams{
		{ID: currPrimary.ID, IsPrimary: &currPrimary.IsPrimary},
		{ID: nextPrimary.ID, IsPrimary: &nextPrimary.IsPrimary},
	}

	for _, update := range emailsToUpdate {
		err = s.emailRepo.Update(ctx, update)
		if err != nil {
			msg := fmt.Sprintf("error on updating email with id %s", update.ID)
			s.log.Errorw(ctx, msg, logger.Err(err))
			return NewBadRequestError(msg, err)
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
