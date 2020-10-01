package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/tjblackheart/shorty/db"
)

type (
	// Config holds the application configuration
	Config struct {
		Port        string
		DQN         string
		Secret      string
		Credentials Creds
	}

	// App is an application instance
	App struct {
		cfg         *Config
		db          db.Repository
		session     *scs.SessionManager
		templates   string
		credentials Creds
		manifest    Manifest
	}

	// Flash holds a flash message.
	Flash struct{ Type, Message string }

	// Creds holds admin credentials, set by .env
	Creds struct{ User, BcryptPass string }

	// Manifest holds a list of compiled webpack assets.
	Manifest map[string]string
)
