package db

import (
	"ecommerce/config"
	"ecommerce/utils"
)

type AuthRepoPG struct {
	PG *config.PG
}

func NewAuthRepoPG(pg *config.PG) AuthRepoPG {
	return AuthRepoPG{
		PG: pg,
	}
}

func (ar AuthRepoPG) Login(username string, password string) (bool, error) {
	var (
		id           int
		hashPassword string
		err          error
	)

	hashPassword, err = utils.HashPassword(password)

	if err != nil {
		return false, err
	}

	err = ar.PG.DB.QueryRow("SELECT id FROM users WHERE username = $1 OR email = $2 AND password = $3", username, username, hashPassword).Scan(&id)

	if err != nil {
		return false, err
	}

	if id == 0 {
		return false, nil
	}

	return true, nil
}

func (ar AuthRepoPG) Register(name string, username, email, password string) error {
	var (
		hashPassword string
		exist        bool
		err          error
	)

	err = ar.PG.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)", username, email).Scan(&exist)

	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	hashPassword, err = utils.HashPassword(password)

	if err != nil {
		return err
	}

	_, err = ar.PG.DB.Exec("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)", name, username, email, hashPassword)

	if err != nil {
		return err
	}

	return nil
}
