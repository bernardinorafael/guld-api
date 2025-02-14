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

	existingUser, err := s.userRepo.FindByID(ctx, dto.UserID)
	if err != nil {
		msg := "error on retrieve user by id"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	if existingUser == nil {
		msg := "user not found"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewNotFoundError(msg, err)
	}

	existingEmail, err := s.emailRepo.FindByEmail(ctx, dto.Email)
	if err != nil {
		msg := "error on retrieve email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}
	if existingEmail != nil {
		msg := "email already exists"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewConflictError(msg, ResourceAlreadyTaken, err, []Field{
			{Field: "email", Msg: msg},
		})
	}

	if !dto.IsPrimary {
		err = s.emailRepo.Insert(ctx, email.AdditionalEmail{
			ID:         util.GenID("email"),
			UserID:     dto.UserID,
			Email:      dto.Email,
			IsPrimary:  dto.IsPrimary,
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

	err = s.userRepo.Update(ctx, user.PartialEntity{
		ID:           existingUser.User.ID,
		EmailAddress: &dto.Email,
	})
	if err != nil {
		msg := "error on update user email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	primaryEmail, err := s.emailRepo.FindPrimary(ctx, dto.UserID)
	if err != nil {
		msg := "error on find primary email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	primaryEmail.IsPrimary = false
	err = s.emailRepo.Update(ctx, email.EmailUpdateParams{
		ID:        primaryEmail.ID,
		IsPrimary: &primaryEmail.IsPrimary,
	})
	if err != nil {
		msg := "error on update primary email"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return NewBadRequestError(msg, err)
	}

	err = s.emailRepo.Insert(ctx, email.AdditionalEmail{
		ID:         util.GenID("email"),
		UserID:     dto.UserID,
		Email:      dto.Email,
		IsPrimary:  dto.IsPrimary,
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
