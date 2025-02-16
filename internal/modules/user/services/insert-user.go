package usersvc

import (
	"context"
	"errors"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/lib/pq"
)

func (s svc) Create(ctx context.Context, dto user.UserRegisterParams) (userId string, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	newUser, err := user.NewUser(dto.FullName, dto.Username, dto.PhoneNumber, dto.EmailAddress)
	if err != nil {
		// TODO: add fields to error
		msg := "error on init user entity"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", NewValidationFieldError(msg, err, nil)
	}

	if err := s.userRepo.Create(ctx, newUser.Store()); err != nil {
		msg := "failed to create user"
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			var appErr = NewConflictError(msg, ResourceAlreadyTaken, err, nil)

			field := util.ExtractFieldFromDetail(pqErr.Detail)
			s.log.Errorw(ctx, msg, logger.Err(err))
			appErr.AddField(field, field+" already exists")
			return "", appErr
		}
		return "", NewBadRequestError(msg, err)
	}

	return newUser.ID(), nil
}
