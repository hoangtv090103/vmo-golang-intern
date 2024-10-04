package db

import (
	"ecommerce/internal/auth/domain"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Login(username, password string) (domain.Auth, error) {
	args := m.Called(username, password)
	return args.Get(0).(domain.Auth), args.Error(1)
}

func (m *MockAuthRepository) Register(auth domain.Auth) error {
	args := m.Called(auth)
	return args.Error(0)
}
