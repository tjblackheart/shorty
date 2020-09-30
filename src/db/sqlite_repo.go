package db

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite driver only needed here.
	"github.com/tjblackheart/shorty/models"
)

type sqliteRepo struct {
	db *sql.DB
}

func Connect(path string) (Repository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	repo := &sqliteRepo{db}
	if err := repo.init(); err != nil {
		return nil, err
	}

	log.Println("Database connected.")
	return repo, nil
}

func (r sqliteRepo) Disconnect() {
	r.db.Close()
}

func (r sqliteRepo) init() error {
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

func (r sqliteRepo) Find() ([]*models.Shorty, error) {
	query := "SELECT * FROM shorty ORDER BY created DESC"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shorties := []*models.Shorty{}
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

func (r sqliteRepo) FindOneByID(id int) (*models.Shorty, error) {
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

func (r sqliteRepo) FindOneByShortLink(shortLink string) (*models.Shorty, error) {
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

func (r sqliteRepo) FindOneByURL(url string) (*models.Shorty, error) {
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

func (r sqliteRepo) DeleteOne(shortLink string) error {
	query := "DELETE FROM shorty WHERE short_link = ?"
	if _, err := r.db.Exec(query, shortLink); err != nil {
		return err
	}

	return nil
}

func (r sqliteRepo) DeleteMany() error {
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

func (r sqliteRepo) Save(s *models.Shorty) error {
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

func (r sqliteRepo) SaveMany(list []*models.Shorty) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, s := range list {
		if err = r.Save(s); err != nil {
			if _, ok := err.(ErrUnique); ok {
				log.Printf("SaveMany: Entry %s already exists.\n", s.Shorty)
				continue
			}

			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r sqliteRepo) Update(s *models.Shorty) error {
	query := "UPDATE shorty SET clicks = ? WHERE id = ?"
	if _, err := r.db.Exec(query, s.Clicks, s.ID); err != nil {
		return err
	}
	return nil
}
