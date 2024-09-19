package usecase

import (
	"ecommerce/internal/user/domain"
	"ecommerce/internal/user/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) AddUser(user domain.User) error {
	return uc.userRepo.Create(user)
}

func (uc *UserUseCase) GetAllUsers() ([]domain.User, error) {
	return uc.userRepo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id int) (domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *UserUseCase) GetUserByUsername(username string) (domain.User, error) {
	return uc.userRepo.GetByUsername(username)
}

func (uc *UserUseCase) UpdateUser(user domain.User) error {
	return uc.userRepo.Update(user)
}

func (uc *UserUseCase) DeleteUser(id int) error {
	return uc.userRepo.Delete(id)
}
