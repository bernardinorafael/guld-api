package usersvc

import (
	"context"
	"fmt"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/phone"
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

func (s svc) SetPrimaryPhone(ctx context.Context, userId, phoneId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	phones, err := s.phoneRepo.FindAllByUser(ctx, userId)
	if err != nil {
		msg := "error on find all phones"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	var currPrimary *phone.AdditionalPhone
	var nextPrimary *phone.AdditionalPhone

	for _, v := range phones {
		if v.IsPrimary {
			currPrimary = &v
		}
		if v.ID == phoneId {
			nextPrimary = &v
		}
	}

	if nextPrimary == nil {
		msg := fmt.Sprintf("not found phone with id %s", phoneId)
		s.log.Errorw(ctx, msg, nil)
		return NewNotFoundError(msg, nil)

	}

	currPrimary.IsPrimary = false
	nextPrimary.IsPrimary = true

	phonesToUpdate := []phone.PhoneUpdateParams{
		{ID: currPrimary.ID, IsPrimary: &currPrimary.IsPrimary},
		{ID: nextPrimary.ID, IsPrimary: &nextPrimary.IsPrimary},
	}

	// This isn't the best approach
	// in the future we should use a batch update
	for _, update := range phonesToUpdate {
		err = s.phoneRepo.Update(ctx, update)
		if err != nil {
			msg := fmt.Sprintf("error on updating phone with id %s", update.ID)
			s.log.Errorw(ctx, msg, logger.Err(err))
			return NewBadRequestError(msg, err)
		}
	}

	err = s.userRepo.Update(ctx, user.PartialEntity{
		ID:          userId,
		PhoneNumber: &nextPrimary.Phone,
	})
	if err != nil {
		msg := "error on updating user phone number"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
