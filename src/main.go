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
		Debug:    false,
		Port:     os.Getenv("APP_PORT"),
		Secret:   os.Getenv("APP_SECRET"),
		DSN:      "/data/db.sqlite",
		ViewsDir: "./templates",
		Credentials: app.Creds{
			User:       os.Getenv("APP_USER"),
			BcryptPass: os.Getenv("APP_BCRYPT_PW"),
		},
	}

	if os.Getenv("APP_DEBUG") == "true" {
		cfg.Debug = true
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
