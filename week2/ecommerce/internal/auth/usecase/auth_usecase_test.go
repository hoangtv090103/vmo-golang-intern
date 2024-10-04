package usecase

import (
	"ecommerce/internal/auth/domain"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Login(username, password string) (domain.Auth, error) {
	args := m.Called(username, password)
	return args.Get(0).(domain.Auth), args.Error(1)
}

func (m *MockAuthUsecase) Register(auth domain.Auth) error {
	args := m.Called(auth)
	return args.Error(0)
}

func TestAuthUsecase_Login(t *testing.T) {
	mockRepo := new(MockAuthUsecase)
	usecase := NewAuthUsecase(mockRepo)

	auth := domain.Auth{
		Username: "testuser",
		Password: "testpassword",
	}

	mockRepo.On("Login", auth.Username, auth.Password).Return(domain.Auth{}, nil)
	token, err := usecase.Login(auth.Username, auth.Password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Error(t *testing.T) {
	mockRepo := new(MockAuthUsecase)
	usecase := NewAuthUsecase(mockRepo)

	mockRepo.On("Login", "wronguser", "wrongpassword").Return(domain.Auth{}, errors.New("invalid credentials"))

	result, err := usecase.Login("wronguser", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register(t *testing.T) {
	mockRepo := new(MockAuthUsecase)
	usecase := NewAuthUsecase(mockRepo)

	auth := domain.Auth{
		Name:     "Test User",
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "testpassword",
	}

	mockRepo.On("Register", auth).Return(nil)
	err := usecase.Register(auth)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
