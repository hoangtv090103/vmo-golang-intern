package usecase

import (
	"ecommerce/internal/auth/domain"
	"ecommerce/internal/auth/repository"
	"ecommerce/internal/auth/utils"
)

type AuthUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(authRepo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func (uc *AuthUsecase) Login(username string, password string) (string, error) {
	auth, err := uc.authRepo.Login(username, password)

	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(auth.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}

//func (uc *AuthUsecase) ForgetPassword(username string) error {
//	return uc.authRepo.ForgetPassword(username)
//}

func (uc *AuthUsecase) Register(auth domain.Auth) error {
	return uc.authRepo.Register(auth)
}
