package domain

type Auth struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Auth) GetID() int {
	return a.ID
}

func (a *Auth) GetUsername() string {
	return a.Username
}

func (a *Auth) GetEmail() string {
	return a.Email
}

func (a *Auth) GetPassword() string {
	return a.Password
}

func (a *Auth) SetUsername(username string) {
	a.Username = username
}

func (a *Auth) SetPassword(password string) {
	a.Password = password
}

func (a *Auth) SetEmail(email string) {
	a.Email = email
}
