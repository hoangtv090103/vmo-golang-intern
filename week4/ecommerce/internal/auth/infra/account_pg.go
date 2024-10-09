package infra

import (
	"context"
	"database/sql"
	"ecommerce/internal/auth/entity"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AccountPGRepository struct {
	DB *sql.DB
}

func NewAccountPGRepository(db *sql.DB) *AccountPGRepository {
	return &AccountPGRepository{DB: db}
}

func (r *AccountPGRepository) Login(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	query := `
        SELECT id, user_id, role_id, name, username, email, password, created_at, updated_at
        FROM accounts
        WHERE username = ? OR email = ?
    `
	// Rebind the query to use $1, $2, etc.
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	foundAccount := &entity.Account{}
	err := r.DB.QueryRowContext(ctx, query, account.Username, account.Email).Scan(foundAccount)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("account not found") // Account not found
	} else if err != nil {
		return nil, err
	}

	return foundAccount, nil
}

func (r *AccountPGRepository) Register(ctx context.Context, account *entity.Account) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	// Create User
	userQuery := `
			INSERT INTO users (name, username, email)
			VALUES (?, ?, ?)
			RETURNING id
		`
	userQuery = sqlx.Rebind(sqlx.DOLLAR, userQuery)

	err = tx.QueryRowContext(
		ctx,
		userQuery,
		account.Name,
		account.Username,
		account.Email,
	).Scan(&account.UserID)

	query := `
            INSERT INTO accounts (user_id, username, email, password)
            VALUES (?, ?, ?, ?)
            RETURNING id
        `
	// Rebind the query to use $1, $2, etc.
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err = tx.QueryRowContext(
		ctx,
		query,
		account.UserID,
		account.Username,
		account.Email,
		account.Password,
	).Scan(&account.ID)

	if err != nil {
		// Check for unique constraint violation
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				if strings.Contains(pqErr.Message, "username") {
					return errors.New("username already exists")
				}
				if strings.Contains(pqErr.Message, "email") {
					return errors.New("email already exists")
				}
			}
		}
		return err
	}

	// Add user role
	roleQuery := `
			INSERT INTO user_roles (user_id, SELECT id FROM roles WHERE role_name = ?)
			VALUES (?, ?)
		`
	roleQuery = sqlx.Rebind(sqlx.DOLLAR, roleQuery)

	_, err = tx.ExecContext(
		ctx,
		roleQuery,
		"user",
		account.UserID,
	)

	if err != nil {
		return err
	}

	err = tx.Commit()

	return nil
}

func (r *AccountPGRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	query := "SELECT a.id, a.user_id, a.username, a.email, a.password, a.created_at, a.updated_at, COALESCE(r.role_name, '') FROM accounts a LEFT JOIN user_roles ur ON ur.auth_id = a.id LEFT JOIN roles r ON ur.role_id = r.id WHERE a.username = ?"
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	row := r.DB.QueryRowContext(ctx, query, username)

	account := &entity.Account{}
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Username,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return account, nil
}
