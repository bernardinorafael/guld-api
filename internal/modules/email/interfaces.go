package email

import (
	"context"
)

type RepositoryInterface interface {
	Create(ctx context.Context, email AdditionalEmail) error
	Update(ctx context.Context, email EmailUpdateParams) error
	GetByID(ctx context.Context, emailId string) (*AdditionalEmail, error)
	GetByEmail(ctx context.Context, email string) (*AdditionalEmail, error)
	GetAllByUser(ctx context.Context, userId string) ([]AdditionalEmail, error)
	Delete(ctx context.Context, userId, emailId string) error
	GetPrimary(ctx context.Context, userId string) (*AdditionalEmail, error)
}
