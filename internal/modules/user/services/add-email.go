package usersvc

import (
	"context"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) AddEmail(ctx context.Context, dto user.CreateEmailParams) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	allEmails, err := s.emailRepo.FindAllByUser(ctx, dto.UserID)
	if err != nil {
		msg := "error on find all emails by user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	// TODO: If needed, transform into a table along with other settings
	if len(allEmails) >= maxEmailAndPhoneByUser {
		msg := "user already has the maximum number of emails"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewForbiddenError(msg, MaxLimitResourceReached, err)
	}

	for _, v := range allEmails {
		if v.Email == dto.Email {
			msg := "email already exists"
			s.log.Errorw(ctx, msg, logger.Err(err))
			return NewConflictError(msg, ResourceAlreadyTaken, err, []Field{
				{Field: "email", Msg: msg},
			})
		}
	}

	err = s.emailRepo.Insert(ctx, email.AdditionalEmail{
		ID:         util.GenID("email"),
		UserID:     dto.UserID,
		Email:      dto.Email,
		IsPrimary:  false,
		IsVerified: false,
		Created:    time.Now(),
		Updated:    time.Now(),
	})
	if err != nil {
		msg := "error on insert email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	return nil
}
