package usecase

import "ecommerce/internal/auth/repository"

type AuthUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(authRepo repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func (uc *AuthUsecase) Login(username string, password string) (bool, error) {
	return uc.authRepo.Login(username, password)
}

func (uc *AuthUsecase) ForgetPassword(username string) error {
	return uc.authRepo.ForgetPassword(username)
}

func (uc *AuthUsecase) Register(username string, password string) error {
	return uc.authRepo.Register(username, password)
}
