package handler

import (
	"bytes"
	"ecommerce/internal/auth/domain"
	"ecommerce/internal/auth/usecase"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAuthUsecase struct {
	mock.Mock
	usecase.AuthUsecase
}

func (m *MockAuthUsecase) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUsecase) Register(auth domain.Auth) error {
	args := m.Called(auth)
	return args.Error(0)
}

func TestAuthHandler_Login(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := NewAuthHandler(&mockUsecase.AuthUsecase)

	auth := domain.Auth{
		Username: "hoangtv",
		Password: "1",
	}

	token := "mockToken"
	mockUsecase.On("Login", auth.Username, auth.Password).Return(token, nil)

	body, _ := json.Marshal(auth)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, token, rr.Body.String())
	mockUsecase.AssertExpectations(t)
}
