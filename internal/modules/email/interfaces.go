package email

import (
	"context"
)

type RepositoryInterface interface {
	Insert(ctx context.Context, email AdditionalEmail) error
	Update(ctx context.Context, email EmailUpdateParams) error
	FindByID(ctx context.Context, emailId string) (*AdditionalEmail, error)
	FindByEmail(ctx context.Context, email string) (*AdditionalEmail, error)
	FindAllByUser(ctx context.Context, userId string) ([]AdditionalEmail, error)
	Delete(ctx context.Context, userId, emailId string) error
}
