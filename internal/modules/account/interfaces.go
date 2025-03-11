package account

import (
	"context"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, acc EntityWithUser) error
	FindByID(ctx context.Context, accountId string) (*EntityWithUser, error)
	FindByUserID(ctx context.Context, userId string) (*Entity, error)
	FindByUsername(ctx context.Context, username string) (*EntityWithUser, error)
	Update(ctx context.Context, acc Entity) error
}

type ServiceInterface interface {
	Login(ctx context.Context, username, password string) (*AccountPayload, error)
	GetSignedInAccount(ctx context.Context) (*EntityWithUser, error)
	ChangePassword(ctx context.Context, userId string, old string, new string) error
}
