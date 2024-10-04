package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	token, err := GenerateJWT("testuser", "user")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	token, err := GenerateJWT("testuser", "user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	_, err := ValidateJWT("invalidtoken")
	assert.Error(t, err)
}
