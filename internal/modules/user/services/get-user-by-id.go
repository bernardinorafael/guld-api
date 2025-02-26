package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) FindByID(ctx context.Context, userId string) (*user.CompleteEntity, error) {
	response, err := s.userRepo.FindCompleteByID(ctx, userId)
	if err != nil {
		s.log.Errorw(ctx, "error on retrieve user", logger.Err(err))
		return nil, NewNotFoundError("error on retrieve user", err)
	}

	codes, err := s.emailService.FindActiveCodesByUser(ctx, userId)
	if err != nil {
		// Thats not a critical error in this context, so we can continue
		s.log.Errorw(ctx, "error on retrieve active codes", logger.Err(err))
	}

	meta := make(map[string]any)

	if len(codes) > 0 {
		activeCodes := make([]email.ValidationEntity, 0)
		for _, code := range codes {
			if code.IsValid {
				activeCodes = append(activeCodes, code)
			}
		}

		emailIDs := make([]string, 0)
		for _, code := range activeCodes {
			emailIDs = append(emailIDs, code.EmailID)
		}

		meta["pending_email_confirmations"] = emailIDs
	}

	return &user.CompleteEntity{
		User:   response.User,
		Emails: response.Emails,
		Meta:   meta,
	}, nil
}
