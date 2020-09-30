package app

import (
	"encoding/gob"
	"math/rand"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/db"
)

// Create builds a new App
func Create(cfg *Config) *App {
	rand.Seed(time.Now().UnixNano())

	repo, err := db.Connect(cfg.DQN)
	if err != nil {
		log.Fatalln(err)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	gob.Register(Flash{})

	app := &App{
		cfg:         cfg,
		repo:        repo,
		session:     session,
		templates:   "templates",
		credentials: cfg.Credentials,
	}

	app.parseManifest()

	return app
}

// Serve starts the server
func (app App) Serve() {
	defer app.repo.Disconnect()

	srv := http.Server{
		Addr:         app.cfg.Addr,
		Handler:      app.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Infof("App listening at %s ...\n", app.cfg.Addr)
	log.Fatalln(srv.ListenAndServe())
}

//

func (app App) err(pkg, msg string) {
	log.Errorf("%s: %s", pkg, msg)
}
