package app

import (
	"encoding/gob"
	"math/rand"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/store"
)

// Create builds a new App
func Create(cfg *Config) *App {
	rand.Seed(time.Now().UnixNano())

	store, err := store.SQLite(cfg.DSN)
	if err != nil {
		log.Fatalln(err)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	gob.Register(flash{})

	app := &App{
		cfg:     cfg,
		store:   store,
		session: session,
		policy:  bluemonday.UGCPolicy(),
	}

	app.initTemplates()

	return app
}

// Serve starts the server
func (app App) Serve() {
	defer app.store.CloseDB()

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
