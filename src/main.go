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
		Port:   os.Getenv("APP_PORT"),
		Secret: os.Getenv("APP_SECRET"),
		DQN:    "/data/db.sqlite",
		Credentials: app.Creds{
			User:       os.Getenv("APP_USER"),
			BcryptPass: os.Getenv("APP_BCRYPT_PW"),
		},
	}

	if cfg.Port == "" {
		cfg.Port = ":3000"
		log.Infof("APP_PORT not set, defaulting to %s", cfg.Port)
	}

	if string(cfg.Port[0]) != ":" {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.Credentials.User == "" || cfg.Credentials.BcryptPass == "" {
		return nil, errors.New("missing credentials")
	}

	if cfg.Secret == "" {
		return nil, errors.New("missing secret")
	}

	return cfg, nil
}
