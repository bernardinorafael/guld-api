package usersvc

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/phone"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) AddPhone(ctx context.Context, dto phone.CreatePhoneParams) error {
	allPhones, err := s.phoneRepo.FindAllByUser(ctx, dto.UserID)
	if err != nil {
		msg := "error on find all phones by user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	for _, v := range allPhones {
		if v.Phone == dto.Phone {
			msg := "phone already exists"
			s.log.Errorw(ctx, msg, logger.Err(err))
			return NewConflictError(msg, ResourceAlreadyTaken, err, []Field{
				{Field: "phone", Msg: msg},
			})
		}
	}

	// TODO: If needed, transform into a table along with other settings
	if len(allPhones) >= maxEmailAndPhoneByUser {
		msg := "user already has the maximum number of phones"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewForbiddenError(msg, MaxLimitResourceReached, err)
	}

	err = s.phoneRepo.Insert(ctx, phone.AdditionalPhone{
		ID:         util.GenID("phone"),
		UserID:     dto.UserID,
		Phone:      dto.Phone,
		IsPrimary:  dto.IsPrimary,
		IsVerified: false,
		Created:    time.Now(),
		Updated:    time.Now(),
	})
	if err != nil {
		msg := "error on insert phone"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	return nil
}
