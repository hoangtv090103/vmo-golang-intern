package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth_Getters(t *testing.T) {
	auth := Auth{
		ID:       1,
		UserID:   2,
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "password123",
	}

	assert.Equal(t, 1, auth.GetID())
	assert.Equal(t, 2, auth.UserID)
	assert.Equal(t, "johndoe", auth.GetUsername())
	assert.Equal(t, "john@example.com", auth.GetEmail())
	assert.Equal(t, "password123", auth.GetPassword())
}

func TestAuth_Setters(t *testing.T) {
	auth := Auth{}

	auth.SetUsername("janedoe")
	auth.SetPassword("newpassword")
	auth.SetEmail("jane@example.com")

	assert.Equal(t, "janedoe", auth.GetUsername())
	assert.Equal(t, "newpassword", auth.GetPassword())
	assert.Equal(t, "jane@example.com", auth.GetEmail())
}
