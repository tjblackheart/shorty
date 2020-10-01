package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite driver
	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/models"
)

type repository struct {
	db *sql.DB
}

// SQLite connects a sqlite3 DB for the given file path
func SQLite(path string) (Repository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	r := &repository{db}
	if err := r.init(); err != nil {
		return nil, err
	}

	log.Info("Database connected.")
	return r, nil
}

func (r repository) Close() {
	r.db.Close()
}

func (r repository) init() error {
	query := `CREATE TABLE IF NOT EXISTS shorty (
		id INTEGER NOT NULL PRIMARY KEY,
		link TEXT NOT NULL,
		short_link VARCHAR(6) UNIQUE NOT NULL,
		clicks INT DEFAULT 0,
		created DATETIME NOT NULL,
		ip VARCHAR(100) NOT NULL
	)`

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}

//

func (r repository) Find() ([]*models.Shorty, error) {
	query := "SELECT * FROM shorty ORDER BY created DESC"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shorties []*models.Shorty
	for rows.Next() {
		var s models.Shorty
		err = rows.Scan(&s.ID, &s.URL, &s.Shorty, &s.Clicks, &s.CreatedAt, &s.IP)
		if err != nil {
			return nil, err
		}
		shorties = append(shorties, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shorties, nil
}

func (r repository) FindOneByID(id int) (*models.Shorty, error) {
	var s models.Shorty

	query := "SELECT * FROM shorty WHERE id = ? LIMIT 1"
	if err := r.db.QueryRow(query, id).Scan(&s.ID, &s.URL, &s.Shorty, &s.Clicks, &s.CreatedAt, &s.IP); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &s, nil
}

func (r repository) FindOneByShortLink(shortLink string) (*models.Shorty, error) {
	var s models.Shorty

	query := "SELECT * FROM shorty WHERE short_link = ? LIMIT 1"
	if err := r.db.QueryRow(query, shortLink).Scan(&s.ID, &s.URL, &s.Shorty, &s.Clicks, &s.CreatedAt, &s.IP); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &s, nil
}

func (r repository) FindOneByURL(url string) (*models.Shorty, error) {
	var s models.Shorty

	query := "SELECT * FROM shorty WHERE link = ? LIMIT 1"
	if err := r.db.QueryRow(query, url).Scan(&s.ID, &s.URL, &s.Shorty, &s.Clicks, &s.CreatedAt, &s.IP); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound{}
		}
		return nil, err
	}

	return &s, nil
}

func (r repository) DeleteOne(shortLink string) error {
	query := "DELETE FROM shorty WHERE short_link = ?"
	if _, err := r.db.Exec(query, shortLink); err != nil {
		return err
	}

	return nil
}

func (r repository) DeleteMany() error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
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

func (r repository) Save(s *models.Shorty) error {
	created := s.CreatedAt
	if s.CreatedAt.IsZero() {
		created = time.Now()
	}

	query := "INSERT INTO shorty (link, short_link, clicks, created, ip) VALUES (?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query, s.URL, s.Shorty, s.Clicks, created, s.IP)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return ErrUnique{}
		}
		return err
	}

	return nil
}

func (r repository) SaveMany(list []*models.Shorty) (int, error) {
	count := 0
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return count, err
	}

	for _, s := range list {
		if err = r.Save(s); err != nil {
			if _, ok := err.(ErrUnique); ok {
				log.Infof("SaveMany: Entry %s already exists, skipping.", s.Shorty)
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

func (r repository) Update(s *models.Shorty) error {
	query := "UPDATE shorty SET clicks = ? WHERE id = ?"
	if _, err := r.db.Exec(query, s.Clicks, s.ID); err != nil {
		return err
	}
	return nil
}
