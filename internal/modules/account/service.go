package account

import (
	"context"

	. "github.com/bernardinorafael/internal/_shared/errors"
	"github.com/bernardinorafael/internal/infra/http/middleware"
	"github.com/bernardinorafael/internal/infra/token"
	"github.com/bernardinorafael/internal/mailer"
	"github.com/bernardinorafael/internal/modules/account/session"
	"github.com/bernardinorafael/internal/modules/user"
	"github.com/bernardinorafael/pkg/logger"
)

type svc struct {
	ctx         context.Context
	log         logger.Logger
	repo        RepositoryInterface
	userRepo    user.RepositoryInterface
	sessionRepo session.RepositoryInterface
	mailer      mailer.Mailer
	secretKey   string
}

func NewService(
	ctx context.Context,
	log logger.Logger,
	repo RepositoryInterface,
	userRepo user.RepositoryInterface,
	sessionRepo session.RepositoryInterface,
	mailer mailer.Mailer,
	secretKey string,
) ServiceInterface {
	return &svc{
		ctx:         ctx,
		log:         log,
		repo:        repo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		mailer:      mailer,
		secretKey:   secretKey,
	}
}

func (s svc) GetSignedInAccount(ctx context.Context) (*EntityWithUser, error) {
	claims, ok := ctx.Value(middleware.AuthKey{}).(*token.AccountClaims)
	if !ok {
		return nil, NewBadRequestError("user ID not found in context", nil)
	}

	acc, err := s.repo.FindByID(ctx, claims.AccountID)
	if err != nil {
		return nil, NewBadRequestError("error on get account by id", err)
	}
	if acc == nil {
		return nil, NewNotFoundError("account not found", nil)
	}

	return acc, nil
}
