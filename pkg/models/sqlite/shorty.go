package sqlite

import (
	"database/sql"
	"strings"

	"github.com/tjblackheart/shorty/pkg/models"
)

// ShortyModel wraps a db connection
type ShortyModel struct {
	DB *sql.DB
}

// InitTable creates DB table
func (m *ShortyModel) InitTable() error {
	query := `CREATE TABLE IF NOT EXISTS shorty (
		id INTEGER NOT NULL PRIMARY KEY,
		link TEXT NOT NULL,
		short_link VARCHAR(6) UNIQUE NOT NULL,
		clicks INT DEFAULT 0,
		created DATETIME NOT NULL,
		ip VARCHAR(100) NOT NULL
	)`

	_, err := m.DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// FindAll returns all shorties
func (m *ShortyModel) FindAll() ([]*models.Shorty, error) {
	query := "SELECT * FROM shorty ORDER BY created DESC"
	rows, err := m.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close() //important

	shorties := []*models.Shorty{}
	for rows.Next() {
		s := &models.Shorty{}
		err = rows.Scan(&s.ID, &s.Link, &s.Shorty, &s.Clicks, &s.Created, &s.IP)
		if err != nil {
			return nil, err
		}

		shorties = append(shorties, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shorties, nil
}

// Find finds one shorty by HashId
func (m *ShortyModel) Find(hashID string) (*models.Shorty, error) {
	s := &models.Shorty{}
	query := "SELECT * FROM shorty WHERE short_link = ?"
	err := m.DB.QueryRow(query, hashID).Scan(&s.ID, &s.Link, &s.Shorty, &s.Clicks, &s.Created, &s.IP)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return s, nil
}

// Insert adds a new shorty
func (m *ShortyModel) Insert(s *models.Shorty) error {
	query := "INSERT INTO shorty (link, short_link, created, ip) VALUES (?, ?, DateTime('now'), ?)"
	_, err := m.DB.Exec(query, s.Link, s.Shorty, s.IP)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return models.ErrUnique
		}

		return err
	}

	return nil
}

// AddClick updates a shorty
func (m *ShortyModel) AddClick(s *models.Shorty) error {
	s.Clicks = s.Clicks + 1

	query := "UPDATE shorty SET clicks = ? WHERE id = ?"
	_, err := m.DB.Exec(query, s.Clicks, s.ID)

	if err != nil {
		return err
	}

	return nil
}

// Remove deletes a shorty
func (m *ShortyModel) Remove(hashID string) error {
	query := "DELETE FROM shorty WHERE short_link = ?"
	_, err := m.DB.Exec(query, hashID)

	if err != nil {
		return err
	}

	return nil
}

// RemoveAll removes all entries
func (m *ShortyModel) RemoveAll() error {
	query := "DELETE FROM shorty"
	_, err := m.DB.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
