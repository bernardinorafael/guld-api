package usersvc

import (
	"context"
	"fmt"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/pkg/logger"

	"github.com/bernardinorafael/internal/modules/phone"
)

func (s svc) FindAllPhones(ctx context.Context, userId string) ([]phone.AdditionalPhone, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		msg := "error on retrieve user by id"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}
	if user == nil {
		msg := "user not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewNotFoundError(msg, err)
	}

	phones, err := s.phoneRepo.FindAllByUser(ctx, userId)
	if err != nil {
		msg := fmt.Sprintf("error on finding phones by user_id %s", userId)
		s.log.Errorw(ctx, msg, logger.Err(err))
		return nil, NewBadRequestError(msg, err)
	}

	return phones, nil
}
