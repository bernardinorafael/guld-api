package account

import (
	"context"
	"errors"
	"time"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/_shared/util"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/crypto"
	"github.com/bernardinorafael/pkg/logger"
	"github.com/lib/pq"
)

type svc struct {
	ctx  context.Context
	log  logger.Logger
	repo RepositoryInterface
	t    *token.Token
}

func NewService(
	ctx context.Context,
	log logger.Logger,
	repo RepositoryInterface,
	secretKey string,
) ServiceInterface {
	return &svc{
		ctx:  ctx,
		log:  log,
		repo: repo,
		t:    token.New(ctx, log, secretKey),
	}
}

func (s svc) GetSignedInAccount(ctx context.Context) (*EntityWithUser, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userId, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		s.log.Errorw(
			ctx,
			"user ID not found in context",
			logger.Err(errors.New("user ID not found in context")),
		)
		return nil, NewConflictError(
			"user ID not found in context",
			InvalidCredentials,
			errors.New("user ID not found in context"),
			nil,
		)
	}

	acc, err := s.repo.FindByID(ctx, userId)
	if err != nil {
		s.log.Errorw(ctx, "error on get signed in account", logger.Err(err))
		return nil, NewConflictError("error on get signed in account", InvalidCredentials, err, nil)
	}

	return acc, nil
}

func (s svc) Login(ctx context.Context, username string, password string) (string, *token.AccountClaims, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	acc, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		s.log.Errorw(ctx, "error on get account by username", logger.Err(err))
		return "", nil, NewConflictError("check username and/or password", InvalidCredentials, err, nil)
	}

	if !crypto.PasswordMatches(password, acc.Password) {
		s.log.Errorw(ctx, "password does not match", logger.Err(err))
		return "", nil, NewConflictError("check username and/or password", InvalidCredentials, err, nil)
	}

	// TODO: Implement retrieve user from account
	t, claims, err := s.t.GenToken(
		acc.ID,
		"-", // Email field from user
		"-", // Username field from user
		60*24*time.Minute,
	)
	if err != nil {
		msg := "error on generate token"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, nil)
	}

	// err = util.
	// 	NewEmailClient("re_Bcc5tfNo_JfhsF7wJBR8dHfR6tsKrpbHj").
	// 	SendEmail(
	// 		util.EmailBody{
	// 			To:      "rafaelferreirab2@gmail.com",
	// 			Subject: "Hello",
	// 			Body:    "account created",
	// 		},
	// 	)
	// if err != nil {
	// 	s.log.Errorw(ctx, "error on send email", logger.Err(err))
	// 	return "", nil, NewBadRequestError("error on send email", err)
	// }

	return t, claims, nil
}

func (s svc) Register(ctx context.Context, dto CreateAccountParams) (string, *token.AccountClaims, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

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

	t, claims, err := s.t.GenToken(
		newAcc.ID(),
		newUser.Email(),
		newUser.Username(),
		60*24*time.Minute,
	)
	if err != nil {
		msg := "error on generate token"
		s.log.Errorw(ctx, msg, logger.Err(err))
		return "", nil, NewBadRequestError(msg, nil)
	}

	return t, claims, nil
}
