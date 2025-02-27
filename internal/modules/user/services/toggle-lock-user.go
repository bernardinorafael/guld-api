package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/modules/user"
)

func (s svc) ToggleLock(ctx context.Context, userId string) error {
	foundUser, err := s.userRepo.FindCompleteByID(ctx, userId)
	if err != nil {
		return NewNotFoundError("failed to retrieve user", err)
	}

	user, err := user.NewFromEntity(foundUser.User)
	if err != nil {
		return NewValidationFieldError("error on init user entity", err, nil)
	}

	err = user.ToggleLock()
	if err != nil {
		return NewBadRequestError("error on toggle lock user status", err)
	}

	err = s.userRepo.Update(ctx, user.Store())
	if err != nil {
		return NewBadRequestError("error on toggle lock user status", err)
	}

	return nil
}
