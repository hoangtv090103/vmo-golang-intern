package repository

import "ecommerce/internal/user/domain"

type UserRepository interface {
	Create(user domain.User) error
	GetAll() ([]domain.User, error)
	GetByID(id int) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	GetByEmail(email string) (domain.User, error)
	Update(user domain.User) error
	Delete(id int) error
}
