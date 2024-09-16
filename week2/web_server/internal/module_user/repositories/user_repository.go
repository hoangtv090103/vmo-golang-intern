package repositories

import "week2-clean-architecture/internal/module_user/domain"

// UserRepository is an interface that defines the methods that a user repositories should implement
type UserRepository interface {
	GetAll() ([]*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}
