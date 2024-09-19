package repository

type AuthRepository interface {
	Login(username string, password string) (bool, error)
	ForgetPassword(username string) error
	Register(username string, password string) error
}