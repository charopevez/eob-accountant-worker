package accounts

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, account Account) (string, error)
	FindByEmail(ctx context.Context, email string) (Account, error)
	FindOne(ctx context.Context, uuid string) (Account, error)
	UpdateCredentials(ctx context.Context, account Account) error
	UpdateBio(ctx context.Context, account Account) error
	Delete(ctx context.Context, uuid string) error
}
