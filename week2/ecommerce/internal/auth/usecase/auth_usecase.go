package usecase

import (
	"ecommerce/internal/auth/domain"
	"ecommerce/internal/auth/repository"
	"ecommerce/internal/auth/utils"
)

type AuthUsecase struct {
	AuthRepo repository.AuthRepository
}

func NewAuthUsecase(authRepo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		AuthRepo: authRepo,
	}
}

// ctx.Context
func (uc *AuthUsecase) Login(username string, password string) (string, error) {
	auth, err := uc.AuthRepo.Login(username, password)

	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(auth.GetUsername(), auth.GetRole())

	if err != nil {
		return "", err

	}

	return token, nil
}

func (uc *AuthUsecase) Register(auth domain.Auth) error {
	return uc.AuthRepo.Register(auth)
}
