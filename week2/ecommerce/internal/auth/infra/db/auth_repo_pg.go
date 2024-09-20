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

	tx, err := ar.PG.DB.Begin()

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

	authQuery := "INSERT INTO auth (user_id, username, email, password) VALUES ($1, $2, $3, $4)"
	_, err = ar.PG.DB.Exec(authQuery, auth.UserID, auth.Username, auth.Email, hashPassword)

	if err != nil {
		return err
	}

	return nil
}
