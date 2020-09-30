package models

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/speps/go-hashids"
)

type (
	Shorty struct {
		ID        int       `json:"-"`
		URL       string    `json:"url" db:"link"`
		Shorty    string    `json:"shortLink" db:"short_link"`
		Clicks    int       `json:"clicks"`
		IP        string    `json:"ip"`
		CreatedAt time.Time `json:"createdAt" db:"created"`
	}

	ErrValidation struct{ err string }
	ErrGenerate   struct{ err string }
)

func (e ErrValidation) Error() string {
	return fmt.Sprintf(e.err)
}

func (e ErrGenerate) Error() string {
	return fmt.Sprintf(e.err)
}

func (s Shorty) Validate() error {
	u, err := url.Parse(strings.TrimSpace(s.URL))
	if err != nil {
		return err
	}

	if u.Scheme == "" || u.Host == "" || u.Scheme != "http" && u.Scheme != "https" {
		return ErrValidation{"invalid URL"}
	}

	return nil
}

func (s *Shorty) Generate() error {
	hd := hashids.NewData()
	hd.Salt = fmt.Sprintf("%f", rand.ExpFloat64()*1e9)
	hd.MinLength = 3

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return ErrGenerate{err.Error()}
	}

	id, err := h.Encode([]int{len(s.URL)})
	if err != nil {
		return ErrGenerate{err.Error()}
	}

	s.Shorty = id
	return nil
}
