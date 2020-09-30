package main

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tjblackheart/shorty/app"
)

func main() {
	cfg, err := getCfg()
	if err != nil {
		log.Fatalln("[INIT] " + err.Error())
	}

	app.Create(cfg).Serve()
}

func getCfg() (*app.Config, error) {
	cfg := &app.Config{
		Addr:   os.Getenv("APP_HOST"),
		Secret: os.Getenv("APP_SECRET"),
		DQN:    "/data/db.sqlite",
		Credentials: app.Creds{
			User:       os.Getenv("APP_USER"),
			BcryptPass: os.Getenv("APP_BCRYPT_PW"),
		},
	}

	if cfg.Credentials.User == "" || cfg.Credentials.BcryptPass == "" {
		return nil, errors.New("missing credentials")
	}

	if cfg.Secret == "" {
		return nil, errors.New("missing secret")
	}

	if cfg.Addr == "" {
		return nil, errors.New("no hostname given")
	}

	return cfg, nil
}
