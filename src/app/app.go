package app

import (
	"encoding/gob"
	"math/rand"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/db"
)

// Create builds a new App
func Create(cfg *Config) *App {
	rand.Seed(time.Now().UnixNano())

	db, err := db.SQLite(cfg.DSN)
	if err != nil {
		log.Fatalln(err)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	gob.Register(Flash{})

	app := &App{
		cfg:         cfg,
		db:          db,
		session:     session,
		templates:   "templates",
		credentials: cfg.Credentials,
		policy:      bluemonday.UGCPolicy(),
	}

	app.initTemplates()

	return app
}

// Serve starts the server
func (app App) Serve() {
	defer app.db.Close()

	srv := http.Server{
		Addr:         app.cfg.Port,
		Handler:      app.router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Infof("App listening at %s ...", app.cfg.Port)
	log.Fatalln(srv.ListenAndServe())
}

//

func (app App) err(pkg, msg string) {
	log.Errorf("%s: %s", pkg, msg)
}
