package account

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/crypto"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/lib/pq"
)

const temporaryTokenDuration = 60 * 24 * time.Minute

type svc struct {
	ctx      context.Context
	log      logger.Logger
	repo     RepositoryInterface
	userRepo user.RepositoryInterface
	t        *token.Token
	mailer   mailer.Mailer
}

func NewService(
	ctx context.Context,
	log logger.Logger,
	repo RepositoryInterface,
	userRepo user.RepositoryInterface,
	mailer mailer.Mailer,
	secretKey string,
) ServiceInterface {
	return &svc{
		ctx:      ctx,
		log:      log,
		repo:     repo,
		userRepo: userRepo,
		mailer:   mailer,
		t:        token.New(ctx, log, secretKey),
	}
}

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

	newAcc, err := NewAccountFromEntity(*account)
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

	toStore := newAcc.Store()

	err = s.repo.Update(ctx, toStore)
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

func (s svc) ActivateAccount(ctx context.Context, userId string) error {
	account, err := s.repo.FindByUserID(ctx, userId)
	if err != nil {
		return NewBadRequestError("error on get account by id", err)
	}

	if account.IsActive {
		s.log.Errorw(ctx, "this account is already activated", logger.Err(err))
		return NewConflictError("error on activating account", ExpiredLink, nil, nil)
	}

	newAcc, err := NewAccountFromEntity(*account)
	if err != nil {
		s.log.Errorw(ctx, "error on create account from entity", logger.Err(err))
	}

	newAcc.Activate()
	toStore := newAcc.Store()

	err = s.repo.Update(ctx, toStore)
	if err != nil {
		s.log.Errorw(ctx, "error on update account", logger.Err(err))
		return NewBadRequestError("error on update account", nil)
	}

	s.log.Infow(ctx, "account activated", logger.Any("account", account))

	return nil
}

func (s svc) GetSignedInAccount(ctx context.Context) (*EntityWithUser, error) {
	accId, ok := ctx.Value(middleware.AccIDKey).(string)
	if !ok {
		return nil, NewConflictError(
			"user ID not found in context",
			InvalidCredentials,
			nil,
			nil,
		)
	}

	acc, err := s.repo.FindByID(ctx, accId)
	if err != nil {
		s.log.Errorw(ctx, "error on get account", logger.Err(err))
		return nil, NewConflictError("error on get account", InvalidCredentials, err, nil)
	}

	return acc, nil
}

func (s svc) Login(ctx context.Context, username string, password string) (string, *token.AccountClaims, error) {
	acc, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", nil, NewConflictError("check username and/or password", InvalidCredentials, err, nil)
	}

	if !acc.IsActive {
		return "", nil, NewConflictError("account is not active", DisabledAccount, err, nil)
	}

	if !crypto.PasswordMatches(password, acc.Password) {
		s.log.Errorw(ctx, "password does not match", logger.Err(err))
		return "", nil, NewConflictError("check username and/or password", InvalidCredentials, err, nil)
	}

	// TODO: Implement retrieve user from account
	t, claims, err := s.t.GenToken(
		token.WithParams{
			AccountID: acc.ID,
			UserID:    acc.User.ID,
			Username:  acc.User.Username,
			Email:     acc.User.EmailAddress,
			OrgID:     &acc.Org.ID,
			Duration:  temporaryTokenDuration,
		},
	)
	if err != nil {
		msg := "error on generate token"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, nil)
	}

	return t, claims, nil
}

func (s svc) Register(ctx context.Context, dto CreateAccountParams) (string, *token.AccountClaims, error) {
	newUser, err := user.NewUser(dto.FullName, dto.Username, dto.PhoneNumber, dto.EmailAddress)
	if err != nil {
		msg := "error on create user"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewValidationFieldError(msg, err, nil)
	}

	encrypted, err := crypto.HashPassword(dto.Password)
	if err != nil {
		msg := "error on hash password"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, err)
	}

	newAcc, err := NewAccount(newUser.ID(), encrypted)
	if err != nil {
		msg := "error on create account"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, err)
	}

	err = s.repo.Insert(ctx, newAcc.StoreWithUser(newUser.Store()))
	if err != nil {
		msg := "failed to create account"
		var pqErr *pq.Error
		// 23505 is the code for unique constraint violation
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			var appErr = NewConflictError(msg, ResourceAlreadyTaken, err, nil)

			field := util.ExtractFieldFromDetail(pqErr.Detail)
			s.log.Errorw(ctx, msg, logger.Err(err))
			appErr.AddField(field, field+" already exists")
			return "", nil, appErr
		}
		return "", nil, NewBadRequestError(msg, err)
	}

	go func() {
		link := fmt.Sprintf("http://localhost:3000/activate/%s", newAcc.UserID())
		params := mailer.SendParams{
			From:    mailer.NotificationSender,
			To:      "rafaelferreirab2@gmail.com",
			Subject: "Activate your account",
			File:    "activate_account.html",
			Data:    map[string]any{"Link": link},
		}
		if err := s.mailer.Send(params); err != nil {
			s.log.Errorw(ctx, "error on send email", logger.Err(err))
		}
	}()

	t, claims, err := s.t.GenToken(
		token.WithParams{
			AccountID: newAcc.ID(),
			UserID:    newUser.ID(),
			Email:     newUser.Email(),
			Username:  newUser.Username(),
			OrgID:     nil,
			Duration:  temporaryTokenDuration,
		},
	)
	if err != nil {
		msg := "error on generate token"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, nil)
	}

	return t, claims, nil
}
