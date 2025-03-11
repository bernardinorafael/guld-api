package account

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/pkg/crypto"
	"github.com/bernardinorafael/pkg/logger"
)

func (s svc) ChangePassword(ctx context.Context, userId string, oldPassword string, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		return NewBadRequestError("error on get user by id", err)
	}
	if user == nil {
		return NewNotFoundError("user not found", nil)
	}

	account, err := s.repo.FindByUserID(ctx, userId)
	if err != nil {
		return NewBadRequestError("error on get account by id", err)
	}
	if account == nil {
		return NewNotFoundError("account not found", nil)
	}

	newAcc, err := NewFromDatabase(*account)
	if err != nil {
		return NewBadRequestError("error on create account entity", nil)
	}

	if !crypto.PasswordMatches(oldPassword, newAcc.password) {
		return NewConflictError("passwords does not matches", InvalidCredentials, nil, nil)
	}

	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return NewBadRequestError("failed to encrypt password", nil)
	}

	err = newAcc.ChangePassword(hashedPassword, user.IgnorePasswordPolicy)
	if err != nil {
		return NewBadRequestError("error on change password", err)
	}

	accountData := newAcc.Store()

	err = s.repo.Update(ctx, accountData)
	if err != nil {
		return NewBadRequestError("error on updating account password", err)
	}

	go func() {
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: "Sua senha foi alterada",
			File:    "change_password.html",
		}
		if err := s.mailer.Send(params); err != nil {
			s.log.Errorw(ctx, "error on send email", logger.Err(err))
		}
	}()

	return nil
}
