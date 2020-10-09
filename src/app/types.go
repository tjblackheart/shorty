package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/flosch/pongo2/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tjblackheart/shorty/db"
)

type (
	// Config holds the application configuration
	Config struct {
		Debug       bool   // if true, templates will be recompiled on each request.
		Port        string // the port to listen on.
		DSN         string // database connection string.
		Secret      string // a random string for csrf generation.
		Credentials Creds  // login credentials.
		ViewsDir    string // path to templates
	}

	// App is an application instance
	App struct {
		cfg         *Config
		db          db.Repository
		session     *scs.SessionManager
		credentials Creds
		manifest    Manifest
		policy      *bluemonday.Policy
	}

	// Flash holds a flash message.
	Flash struct{ Type, Message string }

	// Creds holds admin credentials, set by .env
	Creds struct{ User, BcryptPass string }

	// Manifest holds a list of compiled webpack assets.
	Manifest map[string]string

	// Data holds template data as a pongo2 context.
	Data pongo2.Context
)
