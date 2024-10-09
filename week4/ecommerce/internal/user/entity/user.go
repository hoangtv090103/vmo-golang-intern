package entity

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

func (u *User) SetName(name string) *User {
	u.Name = name
	return u
}

func (u *User) SetUsername(username string) *User {
	u.Username = username
	return u
}

func (u *User) SetEmail(email string) *User {
	u.Email = email
	return u
}

func (u *User) SetBalance(balance float64) *User {
	u.Balance = balance
	return u
}
