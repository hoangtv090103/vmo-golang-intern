package usecase

import (
	"context"
	"ecommerce/internal/user/entity"
	"ecommerce/internal/user/repository"
)

type UserUsecase struct {
	userRepo repository.IUser
}

func NewUserUsecase(userRepo repository.IUser) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *UserUsecase) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *UserUsecase) GetByUserID(ctx context.Context, id int) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *UserUsecase) GetByUserUsername(ctx context.Context, username string) (*entity.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *UserUsecase) GetByUserEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int) error {
	return u.userRepo.Delete(ctx, id)
}
