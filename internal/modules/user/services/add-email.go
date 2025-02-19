package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) AddEmail(ctx context.Context, dto email.CreateParams) error {
	emails, err := s.emailService.FindAllByUser(ctx, dto.UserID)
	if err != nil {
		return err
	}

	// TODO: If needed, transform into a table along with other settings
	if len(emails) >= maxEmailAndPhoneByUser {
		s.log.Errorw(ctx, "max number of emails reached", logger.Err(err))
		return NewForbiddenError("max number of emails reached", MaxLimitResourceReached, err)
	}

	for _, v := range emails {
		if v.Email == dto.Email {
			msg := "email already exists"
			s.log.Errorw(ctx, msg, logger.Err(err))
			return NewConflictError(
				msg,
				ResourceAlreadyTaken,
				err,
				[]Field{{Field: "email", Msg: msg}},
			)
		}
	}

	_, err = s.emailService.Create(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}
