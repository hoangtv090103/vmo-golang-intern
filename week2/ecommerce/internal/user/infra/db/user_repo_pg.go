package db

import (
	"context"
	"database/sql"
	"ecommerce/config"
	"ecommerce/internal/user/domain"
	"strconv"
)

type UserRepoPG struct {
	PG  *config.PG
	RDB *config.Redis
	Ctx context.Context
}

func NewUserRepoPG(pdb *config.PG, rdb *config.Redis) *UserRepoPG {
	return &UserRepoPG{
		PG:  pdb,
		RDB: rdb,
		Ctx: context.Background(),
	}
}

func (u *UserRepoPG) Create(user domain.User) error {
	query := `INSERT INTO users (name, username, email, balance)
              	VALUES ($1, $2, $3, $4, $5) returning id`

	_, err := u.PG.GetDB().Exec(query, user.Name, user.Username, user.Email, user.Balance)

	if err != nil {
		return err
	}

	// // Invalidate any cached user data
	// _ = u.RDB.GetClient().Del(u.Ctx, fmt.Sprintf("user:%d", id)).Err()

	return nil
}

func (u *UserRepoPG) GetAll() ([]domain.User, error) {
	var users []domain.User
	rows, err := u.PG.GetDB().Query("SELECT id, name, username, balance FROM users")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepoPG) GetByID(id int) (domain.User, error) {
	var user domain.User
	err := u.PG.GetDB().QueryRow(
		"SELECT id, name, username, email, balance FROM users WHHERE id = $1",
		id,
	).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Balance)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UserRepoPG) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := u.PG.GetDB().QueryRow(
		"SELECT id, username, email, balance FROM users WHERE username = $1",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Balance)
	if err != nil {
		return domain.User{}, err
	}

	// // Cache the user data
	// err = u.RDB.GetClient().Set(u.Ctx, fmt.Sprintf("user:%d", user.ID), user, 0).Err()

	// if err != nil {
	// 	return domain.User{}, err
	// }

	return user, nil
}

func (u *UserRepoPG) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	err := u.PG.GetDB().QueryRow("SELECT id, username, email, balance FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Balance)
	if err != nil {
		return domain.User{}, err
	}

	// // Cache the user data
	// err = u.RDB.GetClient().Set(u.Ctx, fmt.Sprintf("user:%d", user.ID), user, 0).Err()

	// if err != nil {
	// 	return domain.User{}, err
	// }
	return user, nil
}

func (u *UserRepoPG) Update(user domain.User) error {

	query := "UPDATE users SET"
	params := []interface{}{} // Slice to store the query parameters
	paramCount := 1
	if user.Name != "" {
		query += " name = $" + strconv.Itoa(paramCount) + ","
		params = append(params, user.Name)
		paramCount++
	}
	if user.Username != "" {
		query += " username = $" + strconv.Itoa(paramCount) + ","
		params = append(params, user.Username)
		paramCount++
	}

	if user.Email != "" {
		query += " email = $" + strconv.Itoa(paramCount) + ","
		params = append(params, user.Email)
		paramCount++
	}

	if user.Balance >= 0 {
		query += " balance = $" + strconv.Itoa(paramCount) + ","
		params = append(params, user.Balance)
		paramCount++
	}

	if paramCount == 1 {
		return nil
	}

	// Remove the trailing comma
	query = query[:len(query)-1]

	query += " WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, user.ID)

	_, err := u.PG.GetDB().Exec(query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoPG) Delete(id int) error {
	_, err := u.PG.GetDB().Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	// // Invalidate any cached user data
	// _ = u.RDB.GetClient().Del(u.Ctx, fmt.Sprintf("user:%d", id)).Err()
	return nil
}
