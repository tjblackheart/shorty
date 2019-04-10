package sqlite

import (
	"database/sql"
	"strings"

	"github.com/tjblackheart/shorty/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserModel wraps a db connection
type UserModel struct {
	DB *sql.DB
}

// InitTable creates user table
func (m *UserModel) InitTable() error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		pass VARCHAR(60) NOT NULL,
		created DATETIME NOT NULL
	)`

	_, err := m.DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Create a user
func (m *UserModel) Create(u *models.User) error {
	query := "INSERT INTO user (email, pass, created) VALUES (?, ?, DateTime('now'))"

	_, err := m.DB.Exec(query, u.Email, u.Pass)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrUnique
		}

		return err
	}

	return nil
}

// Find a user by email
func (m *UserModel) Find(email string) (*models.User, error) {
	u := &models.User{}
	query := "SELECT * FROM user WHERE email = ?"

	err := m.DB.QueryRow(query, email).Scan(&u.ID, &u.Email, &u.Pass, &u.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return u, nil
}

// Authenticate a user by email and password
func (m *UserModel) Authenticate(email, password string) (*models.User, error) {
	u, err := m.Find(email)

	if err != nil {
		if err == models.ErrNoRecord {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Pass), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, models.ErrInvalidCredentials
		}
		return nil, err
	}

	return u, nil
}
