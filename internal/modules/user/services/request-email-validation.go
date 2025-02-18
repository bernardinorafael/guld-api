package usersvc

import (
	"context"
	"fmt"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) RequestEmailValidation(ctx context.Context, email, userId string) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	emailId, err := s.getEmail(ctx, email)
	if err != nil {
		return err
	}

	go func() {
		link := fmt.Sprintf("http://localhost:3000/verify-email/%s", emailId)
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: "Ative seu e-mail",
			File:    "email_validation.html",
			Data:    map[string]string{"Link": link},
		}
		if err := s.mailer.Send(params); err != nil {
			s.log.Errorw(ctx, "error sending email", logger.Err(err))
		}
	}()

	return nil
}

func (s svc) getEmail(ctx context.Context, email string) (string, error) {
	addr, err := s.emailRepo.FindByEmail(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "error on get email", logger.Err(err))
		return "", NewBadRequestError("error on get email", err)
	}
	if addr == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return "", NewBadRequestError("email not found", err)
	}

	return addr.ID, nil
}
