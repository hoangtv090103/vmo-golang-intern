package usecases

import (
	"week2-clean-architecture/internal/module_user/domain"
	"week2-clean-architecture/internal/module_user/repositories"
)

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUsecase(userRepo repositories.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) GetUsers() ([]*domain.User, error) {
	return u.userRepo.GetAll()
}

func (u *UserUseCase) GetUserByUsername(username string) (*domain.User, error) {
	return u.userRepo.GetByUsername(username)
}

func (u *UserUseCase) CreateUser(user *domain.User) error {
	return u.userRepo.Create(user)
}
