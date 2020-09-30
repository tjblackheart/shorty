package db

import "github.com/tjblackheart/shorty/models"

type Repository interface {
	Disconnect()
	Find() ([]*models.Shorty, error)
	FindOneByID(id int) (*models.Shorty, error)
	FindOneByShortLink(hashID string) (*models.Shorty, error)
	FindOneByURL(url string) (*models.Shorty, error)
	DeleteOne(shortLink string) error
	DeleteMany() error
	Save(s *models.Shorty) error
	SaveMany(s []*models.Shorty) error
	Update(s *models.Shorty) error
}
