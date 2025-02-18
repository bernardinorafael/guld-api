package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	emails "github.com/bernardinorafael/internal/modules/email"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) FindEmail(ctx context.Context, email string) (*emails.AdditionalEmail, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userId, ok := ctx.Value(middleware.UserIDKey).(string) // ou o tipo apropriado
	if !ok {
		s.log.Errorw(ctx, "user ID not found in context", nil)
		return nil, NewNotFoundError("user ID not found in context", nil)
	}

	foundUser, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		s.log.Errorw(ctx, "error retrieving user", logger.Err(err))
		return nil, NewNotFoundError("error retrieving user", err)
	}
	if foundUser == nil {
		s.log.Errorw(ctx, "user not found", logger.Err(err))
		return nil, NewNotFoundError("user not found", nil)
	}

	foundEmail, err := s.emailRepo.FindByID(ctx, email)
	if err != nil {
		s.log.Errorw(ctx, "error retrieving email", logger.Err(err))
		return nil, NewNotFoundError("error retrieving email", err)
	}
	if foundEmail == nil {
		s.log.Errorw(ctx, "email not found", logger.Err(err))
		return nil, NewNotFoundError("email not found", nil)
	}

	return foundEmail, nil
}
