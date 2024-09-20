package repository

import "ecommerce/internal/auth/domain"

type AuthRepository interface {
	Login(username string, password string) (domain.Auth, error)
	//ForgetPassword(username string) error
	Register(auth domain.Auth) error
}
