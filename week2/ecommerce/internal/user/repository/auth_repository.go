package repository

type AuthRepository interface {
    Login(username string, password string) (bool, error)
    Register(username string, password string) error
    ForgetPassword(username string) error
}