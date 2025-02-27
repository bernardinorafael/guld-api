package usersvc

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"

	"github.com/bernardinorafael/internal/modules/user"
)

func (s svc) UpdateProfile(ctx context.Context, userId string, dto user.UpdateProfileDTO) error {
	foundUser, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return NewNotFoundError("failed to retrieve user", err)
	}

	user, err := user.NewFromEntity(*foundUser)
	if err != nil {
		return NewValidationFieldError("error on init user entity", err, nil)
	}

	err = user.ChangeUsername(dto.Username)
	if err != nil {
		return NewBadRequestError("error on update profile", err)
	}

	err = user.ChangeName(dto.FullName)
	if err != nil {
		return NewBadRequestError("error on update profile", err)
	}

	err = s.userRepo.Update(ctx, user.Store())
	if err != nil {
		return NewBadRequestError("error on update profile", err)
	}

	s.log.Infof(ctx, "user %s updated profile", userId)

	return nil
}
