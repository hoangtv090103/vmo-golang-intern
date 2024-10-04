package repository

import (
	"ecommerce/internal/auth/domain"
	"ecommerce/internal/auth/infra/db"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRepository_Login(t *testing.T) {
	mockRepo := new(db.MockAuthRepository)
	auth := domain.Auth{
		Username: "testuser",
		Password: "hashedpassword",
	}

	mockRepo.On("Login", "testuser", "testpassword").Return(auth, nil)

	result, err := mockRepo.Login("testuser", "testpassword")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestAuthRepository_Login_Error(t *testing.T) {
	mockRepo := new(db.MockAuthRepository)
	mockRepo.On("Login", "wronguser", "wrongpassword").Return(domain.Auth{}, errors.New("invalid credentials"))

	_, err := mockRepo.Login("wronguser", "wrongpassword")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthRepository_Register(t *testing.T) {
	mockRepo := new(db.MockAuthRepository)

	auth := domain.Auth{
		Name:     "Test User",
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "testpassword",
	}

	mockRepo.On("Register", auth).Return(nil)

	err := mockRepo.Register(auth)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
