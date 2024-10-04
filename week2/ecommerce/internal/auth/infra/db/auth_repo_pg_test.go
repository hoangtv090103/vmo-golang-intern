package db

import (
	"ecommerce/config"
	"ecommerce/internal/auth/domain"
	"ecommerce/utils"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthRepoPG_Login(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepoPG(&config.PG{DB: db})
	rows := sqlmock.NewRows([]string{"username", "password"}).
		AddRow("testuser", "hashedpassword")

	mock.ExpectQuery("SELECT username, password FROM users WHERE username = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	auth, err := repo.Login("testuser", "testpassword")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", auth.Username)
	assert.Equal(t, "hashedpassword", auth.Password)
}

func TestAuthRepoPG_Login_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepoPG(&config.PG{
		DB: db,
	})

	mock.ExpectQuery("SELECT username, password FROM users WHERE username = ?").
		WithArgs("wronguser").
		WillReturnError(errors.New("user not found"))

	_, err = repo.Login("wronguser", "wrongpassword")
	assert.Error(t, err)
}

func TestAuthRepoPG_Register(t *testing.T) {
	//Init DB simulator
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepoPG(&config.PG{
		DB: db,
	})

	auth := domain.Auth{
		Name:     "Test User 1",
		Username: "testuser1",
		Email:    "testuser1@gmail.com",
		Password: "testpassword1",
	}

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM users WHERE username = \$1 OR email = \$2\)`).
		WithArgs(auth.Username, auth.Email).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO users \(name, username, email\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(auth.Name, auth.Username, auth.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	var hashedPassword string
	hashedPassword, err = utils.HashPassword(auth.Password)

	if err != nil {
		return
	}

	mock.ExpectExec(`INSERT INTO auth \(user_id, username, email, password\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(1, auth.Username, auth.Email, hashedPassword).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Register(auth)

	assert.NoError(t, err)
}

func TestAuthRepoPG_Register_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewAuthRepoPG(&config.PG{
		DB: db,
	})

	auth := domain.Auth{
		Name:     "Test User",
		Username: "testuser",
		Email:    "testuser@gmail.com",
		Password: "testpassword",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(auth.Name, auth.Username, auth.Email, auth.Password).
		WillReturnError(errors.New("insert error"))

	err = repo.Register(auth)
	assert.Error(t, err)
}
