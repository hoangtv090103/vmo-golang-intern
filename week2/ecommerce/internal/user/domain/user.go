package domain

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Balance  float64 `json:"balance"`
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetBalance() float64 {
	return u.Balance
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetBalance(balance float64) {
	u.Balance = balance
}