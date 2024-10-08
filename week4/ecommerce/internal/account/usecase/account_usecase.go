package usecase

import (
	"context"
	"ecommerce/internal/account/entity"
	"ecommerce/internal/account/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IAccountUsecase interface {
	Register(ctx context.Context, account *entity.Account) error
	Login(ctx context.Context, username, password string) (*entity.Account, error)
	// GetAccountByID(id int) (*entity.Account, error)
	// UpdateAccount(account *entity.Account) error
	// DeleteAccount(id int) error
	// ChangePassword(id int, oldPassword, newPassword string) error
}

type AccountUsecase struct {
	repo repository.AccountRepository
}

func NewAccountUsecase(repo repository.AccountRepository) IAccountUsecase {
	return &AccountUsecase{
		repo: repo,
	}
}

func (u *AccountUsecase) Register(ctx context.Context, account *entity.Account) error {
	// Check if username or email already exist
	existingAccount, err := u.repo.Login(ctx, &entity.Account{Username: account.Username})

	if err == nil && existingAccount != nil {
		return errors.New("username already exists")
	}

	existingAccount, err = u.repo.Login(ctx, &entity.Account{Email: account.Email})
	if err == nil && existingAccount != nil {
		return errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set default role if not provided
	account.
		SetPassword(string(hashedPassword))

	// Save the account
	return u.repo.Register(ctx, account)
}

func (u *AccountUsecase) Login(ctx context.Context, username, password string) (*entity.Account, error) {
	account, err := u.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return account, nil
}
