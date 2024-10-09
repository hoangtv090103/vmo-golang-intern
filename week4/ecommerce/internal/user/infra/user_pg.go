package infra

import (
	"context"
	"database/sql"
	"ecommerce/internal/user/entity"
	"github.com/jmoiron/sqlx"
)

type UserPGRepository struct {
	DB *sql.DB
}

func NewUserPGRepository(db *sql.DB) *UserPGRepository {
	return &UserPGRepository{
		DB: db,
	}
}

func (u *UserPGRepository) Create(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (name, username, email, balance) VALUES (?, ?, ?, ?)"
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	_, err := u.DB.ExecContext(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.Balance,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPGRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	rows, err := u.DB.QueryContext(ctx, "SELECT id, name, username, balance FROM users")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		user := &entity.User{}

		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Balance)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserPGRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	var user *entity.User
	err := u.DB.QueryRowContext(
		ctx,
		"SELECT id, name, username, email, balance FROM users WHERE id = $1",
		id,
	).Scan(user.ID, user.Name, user.Username, user.Email, user.Balance)

	if err != nil {
		return &entity.User{}, err
	}

	return user, nil
}

func (u *UserPGRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	err := u.DB.QueryRowContext(
		ctx,
		"SELECT id, username, email, balance FROM users WHERE username = $1",
		username,
	).Scan(user.ID, user.Username, user.Email, user.Balance)
	if err != nil {
		return &entity.User{}, err
	}

	return user, nil
}

func (u *UserPGRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user *entity.User
	err := u.DB.QueryRowContext(
		ctx,
		"SELECT id, username, email, balance FROM users WHERE email = $1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Balance)
	if err != nil {
		return &entity.User{}, err
	}

	return user, nil
}

func (u *UserPGRepository) Update(ctx context.Context, user *entity.User) error {
	query := "UPDATE users SET"
	params := []interface{}{} // Slice to store the query parameters
	if user.Name != "" {
		query += " name = ?,"
		params = append(params, user.Name)
	}
	if user.Username != "" {
		query += " username = ?,"
		params = append(params, user.Username)
	}

	if user.Email != "" {
		query += " email = ?,"
		params = append(params, user.Email)
	}

	if user.Balance >= 0 {
		query += " balance = ?,"
		params = append(params, user.Balance)
	}

	// Remove the trailing comma
	query = query[:len(query)-1]

	query += " WHERE id = ?"
	params = append(params, user.ID)

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err := u.DB.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPGRepository) Delete(ctx context.Context, id int) error {
	_, err := u.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
