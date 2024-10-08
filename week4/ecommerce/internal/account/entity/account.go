package entity

import "time"

type Account struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type AccountBuilder struct {
	account *Account
}

func NewAccountBuilder() *AccountBuilder {
	return &AccountBuilder{
		account: &Account{},
	}
}

func (a *Account) SetUsername(username string) *Account {
	a.Username = username
	return a
}

func (a *Account) SetPassword(password string) *Account {
	a.Password = password
	return a
}

func (a *Account) SetEmail(email string) *Account {
	a.Email = email
	return a
}

func (a *Account) SetCreatedAt() *Account {
	now := time.Now()
	a.UpdatedAt = &now
	return a
}

func (a *Account) SetUpdatedAt() *Account {
	now := time.Now()
	a.UpdatedAt = &now
	return a
}
