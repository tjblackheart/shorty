package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrUnique             = errors.New("models: UNIQUE key already exists")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)

// Shorty holds a db row.
type Shorty struct {
	ID      int
	Link    string
	Shorty  string `db:"short_link"`
	Clicks  int
	IP      string
	Created time.Time
}

// User holds a db user entry
type User struct {
	ID      int
	Email   string
	Pass    string
	Created time.Time
}
