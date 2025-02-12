package account

import (
	"context"

	"github.com/bernardinorafael/internal/infra/token"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, acc EntityWithUser) error
	FindByID(ctx context.Context, accountId string) (*EntityWithUser, error)
	FindByUserID(ctx context.Context, userId string) (*Entity, error)
	FindByUsername(ctx context.Context, username string) (*EntityWithUser, error)
	Update(ctx context.Context, acc PartialEntity) error
}

type ServiceInterface interface {
	Login(ctx context.Context, username, password string) (string, *token.AccountClaims, error)
	Register(ctx context.Context, dto CreateAccountParams) (string, *token.AccountClaims, error)
	GetSignedInAccount(ctx context.Context) (*EntityWithUser, error)
	ActivateAccount(ctx context.Context, accountId string) error
}
