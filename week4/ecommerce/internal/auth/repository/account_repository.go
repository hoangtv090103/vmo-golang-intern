package repository

import (
	"context"
	"ecommerce/internal/auth/entity"
)

type IAccountRepository interface {
	Login(ctx context.Context, account *entity.Account) (*entity.Account, error)
	Register(ctx context.Context, account *entity.Account) error
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
}
