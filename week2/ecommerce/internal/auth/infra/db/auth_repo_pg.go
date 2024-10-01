package db

import (
	"ecommerce/config"
	"ecommerce/internal/auth/domain"
	"ecommerce/utils"
	"errors"
)

type AuthRepoPG struct {
	PG *config.PG
}

func NewAuthRepoPG(pg *config.PG) AuthRepoPG {
	return AuthRepoPG{
		PG: pg,
	}
}

func (ar AuthRepoPG) Login(username string, password string) (domain.Auth, error) {
	var (
		err  error
		auth domain.Auth
	)

	err = ar.PG.DB.QueryRow(
		"SELECT id, username, email, password FROM auth WHERE username = $1",
		username,
	).Scan(
		&auth.ID,
		&auth.Username,
		&auth.Email,
		&auth.Password,
	)

	if err != nil {
		return domain.Auth{}, err
	}

	// Check if the user exists
	if auth.GetID() == 0 {
		return domain.Auth{}, errors.New("user not found")
	}

	// Verify the password
	if !utils.CheckPasswordHash(password, auth.Password) {
		return domain.Auth{}, errors.New("invalid password")
	}

	// Get user role_id
	err = ar.PG.DB.QueryRow(
		`SELECT role_name FROM roles LEFT JOIN user_roles ON roles.id = user_roles.role_id WHERE user_roles.auth_id = $1`,
		auth.ID,
	).Scan(&auth.Role)

	if err != nil {
		return domain.Auth{}, err
	}

	return auth, nil
}

func (ar AuthRepoPG) Register(auth domain.Auth) error {
	var (
		hashPassword string
		exist        bool
		err          error
	)

	err = ar.PG.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", auth.Username, auth.Email).Scan(&exist)

	if err != nil {
		return err
	}

	if exist {
		return errors.New("user already exists")
	}

	hashPassword, err = utils.HashPassword(auth.Password)

	if err != nil {
		return err
	}

	tx, err := ar.PG.GetDB().Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err := tx.Commit()
		if err != nil {
			return
		}
	}()
	err = ar.PG.DB.QueryRow(
		"INSERT INTO users (name, username, email) VALUES ($1, $2, $3) RETURNING id",
		auth.Name,
		auth.Username,
		auth.Email,
	).Scan(&auth.UserID)

	if err != nil {
		return err
	}

	authQuery := "INSERT INTO auth (user_id, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id"
	err = ar.PG.DB.QueryRow(authQuery, auth.UserID, auth.Username, auth.Email, hashPassword).Scan(&auth.ID)

	if err != nil {
		return err
	}

	userRoleQuery := "INSERT INTO user_roles (auth_id, role_id) VALUES ($1, (SELECT id FROM roles WHERE role_name = 'user'))"
	_, err = ar.PG.DB.Exec(userRoleQuery, auth.ID)

	if err != nil {
		return err
	}

	return nil
}
