package store

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/models"
)

type store struct {
	db *sql.DB
}

// SQLite connects a sqlite3 DB for the given file path
func SQLite(path string) (Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	s := &store{db}
	if err := s.init(); err != nil {
		return nil, err
	}

	log.Info("Database connected.")
	return s, nil
}

func (s store) CloseDB() {
	s.db.Close()
}

func (s store) init() error {
	query := `CREATE TABLE IF NOT EXISTS shorty (
		id INTEGER NOT NULL PRIMARY KEY,
		link TEXT NOT NULL,
		short_link VARCHAR(6) UNIQUE NOT NULL,
		clicks INT DEFAULT 0,
		created DATETIME NOT NULL,
		ip VARCHAR(100) NOT NULL
	)`

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	return nil
}

//

func (s store) Find() ([]*models.Shorty, error) {
	query := "SELECT * FROM shorty ORDER BY created DESC"
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shorties []*models.Shorty
	for rows.Next() {
		var shorty models.Shorty

		if err = rows.Scan(
			&shorty.ID,
			&shorty.URL,
			&shorty.Shorty,
			&shorty.Clicks,
			&shorty.CreatedAt,
			&shorty.IP,
		); err != nil {
			return nil, err
		}

		shorties = append(shorties, &shorty)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shorties, nil
}

func (s store) FindOne(id int) (*models.Shorty, error) {
	var shorty models.Shorty

	query := "SELECT * FROM shorty WHERE id = ? LIMIT 1"
	if err := s.db.QueryRow(query, id).Scan(
		&shorty.ID,
		&shorty.URL,
		&shorty.Shorty,
		&shorty.Clicks,
		&shorty.CreatedAt,
		&shorty.IP,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &shorty, nil
}

func (s store) FindOneByShortLink(shortLink string) (*models.Shorty, error) {
	var shorty models.Shorty

	query := "SELECT * FROM shorty WHERE short_link = ? LIMIT 1"
	if err := s.db.QueryRow(query, shortLink).Scan(
		&shorty.ID,
		&shorty.URL,
		&shorty.Shorty,
		&shorty.Clicks,
		&shorty.CreatedAt,
		&shorty.IP,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &shorty, nil
}

func (s store) FindOneByURL(url string) (*models.Shorty, error) {
	var shorty models.Shorty

	query := "SELECT * FROM shorty WHERE link = ? LIMIT 1"
	if err := s.db.QueryRow(query, url).Scan(
		&shorty.ID,
		&shorty.URL,
		&shorty.Shorty,
		&shorty.Clicks,
		&shorty.CreatedAt,
		&shorty.IP,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &shorty, nil
}

func (s store) DeleteOne(shortLink string) error {
	query := "DELETE FROM shorty WHERE short_link = ?"
	if _, err := s.db.Exec(query, shortLink); err != nil {
		return err
	}

	return nil
}

func (s store) DeleteMany() error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM shorty")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s store) Save(shorty *models.Shorty) error {
	if shorty.CreatedAt.IsZero() {
		shorty.CreatedAt = time.Now()
	}

	query := "INSERT INTO shorty (link, short_link, clicks, created, ip) VALUES (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, shorty.URL, shorty.Shorty, shorty.Clicks, shorty.CreatedAt, shorty.IP)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrUnique{}
		}
		return err
	}

	return nil
}

func (s store) SaveMany(list []*models.Shorty) (int, error) {
	count := 0
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return count, err
	}

	for _, shorty := range list {
		if err = s.Save(shorty); err != nil {
			if _, ok := err.(ErrUnique); ok {
				log.Infof("SaveMany: Entry %s already exists, skipping.", shorty.Shorty)
				continue
			}

			tx.Rollback()
			return count, err
		}
		count++
	}

	if err = tx.Commit(); err != nil {
		return count, err
	}

	return count, nil
}

func (s store) Update(shorty *models.Shorty) error {
	query := "UPDATE shorty SET clicks = ? WHERE id = ?"
	if _, err := s.db.Exec(query, shorty.Clicks, shorty.ID); err != nil {
		return err
	}
	return nil
}
