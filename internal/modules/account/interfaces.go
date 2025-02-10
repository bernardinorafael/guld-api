package account

import (
	"context"

	"github.com/bernardinorafael/internal/infra/token"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, acc EntityWithUser) error
	GetByID(ctx context.Context, accountId string) (*EntityWithUser, error)
	GetByUserID(ctx context.Context, userId string) (*Entity, error)
	GetByUsername(ctx context.Context, username string) (*Entity, error)
}

type ServiceInterface interface {
	Login(ctx context.Context, username, password string) (string, *token.AccountClaims, error)
	Register(ctx context.Context, dto CreateAccountParams) (string, *token.AccountClaims, error)
	GetSignedInAccount(ctx context.Context) (*EntityWithUser, error)
}
