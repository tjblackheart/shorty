package store

import (
	"fmt"

	"github.com/tjblackheart/shorty/models"
)

type (
	Store interface {
		Find() ([]*models.Shorty, error)
		FindOne(id int) (*models.Shorty, error)
		FindOneByShortLink(hashID string) (*models.Shorty, error)
		FindOneByURL(url string) (*models.Shorty, error)
		DeleteOne(shortLink string) error
		DeleteMany() error
		Save(s *models.Shorty) error
		SaveMany(s []*models.Shorty) (int, error)
		Update(s *models.Shorty) error
		CloseDB()
	}

	ErrNotImplemented struct{}
	ErrNotFound       struct{}
	ErrUnique         struct{}
)

func (e ErrNotImplemented) Error() string {
	return fmt.Sprintf("Not implemented.")
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Entry not found.")
}

func (e ErrUnique) Error() string {
	return fmt.Sprintf("Entry already exists")
}
