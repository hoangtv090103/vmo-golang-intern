package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Getter(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "Hoang",
		Username: "hoang",
		Email:    "hoang@gmail.com",
		Balance:  9010.80,
	}

	assert.Equal(t, 1, user.GetID())
	assert.Equal(t, "Hoang", user.GetName())
	assert.Equal(t, "hoang", user.GetUsername())
	assert.Equal(t, "hoang@gmail.com", user.GetEmail())
	assert.Equal(t, 9010.80, user.GetBalance())
}

func TestUser_Setters(t *testing.T) {
	user := User{}

	user.SetUsername("tran")
	user.SetName("Tran")
	user.SetEmail("tran@gmail.com")
	user.SetBalance(1080.0)

	assert.Equal(t, "Tran", user.GetName())
	assert.Equal(t, "tran", user.GetUsername())
	assert.Equal(t, "tran@gmail.com", user.GetEmail())
	assert.Equal(t, 1080.0, user.GetBalance())
}
