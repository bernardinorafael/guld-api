package usersvc

import (
	"context"
	"fmt"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) RequestEmailValidation(ctx context.Context, emailString, userId string) error {
	emailAddr, err := s.emailService.FindByEmail(ctx, emailString)
	if err != nil {
		return err
	}

	if emailAddr.IsVerified {
		s.log.Errorw(ctx, "email already verified")
		return NewForbiddenError("email already verified", BadRequest, nil)
	}

	err = s.emailService.InsertValidation(ctx, email.Validation{
		ID:         util.GenID("val-email"),
		EmailID:    emailAddr.ID,
		IsConsumed: false,
		Created:    time.Now(),
		Expires:    time.Now().Add(time.Minute * 30),
	})
	if err != nil {
		return NewBadRequestError("failed to insert email validaton", nil)
	}

	go func() {
		link := fmt.Sprintf("http://localhost:3000/verify-email/%s", emailAddr.ID)
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: "Ative seu e-mail",
			File:    "email_validation.html",
			Data:    map[string]string{"Link": link},
		}
		if err := s.mailer.Send(params); err != nil {
			s.log.Errorw(
				ctx,
				"error sending email",
				logger.Err(err),
				logger.String("email", params.To),
				logger.String("userId", userId),
			)
			return
		}
		s.log.Infow(ctx, "email validation sent",
			logger.String("email", params.To),
			logger.String("userId", userId),
		)
	}()

	return nil
}
