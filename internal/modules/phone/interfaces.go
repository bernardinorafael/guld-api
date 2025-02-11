package phone

import (
	"context"
)

type RepositoryInterface interface {
	Create(ctx context.Context, phone AdditionalPhone) error
	Update(ctx context.Context, phone PhoneUpdateParams) error
	FindByID(ctx context.Context, phoneId string) (*AdditionalPhone, error)
	FindByPhone(ctx context.Context, phone string) (*AdditionalPhone, error)
	GetAllByUser(ctx context.Context, userId string) ([]AdditionalPhone, error)
	Delete(ctx context.Context, userId, phoneId string) error
	GetPrimary(ctx context.Context, userId string) (*AdditionalPhone, error)
}
