package app

import (
	"github.com/alexedwards/scs/v2"
	"github.com/flosch/pongo2/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tjblackheart/shorty/store"
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
		cfg      *Config
		store    store.Store
		session  *scs.SessionManager
		manifest manifest
		policy   *bluemonday.Policy
	}

	// Creds holds admin credentials, set by .env
	Creds struct{ User, BcryptPass string }

	// Flash holds a flash message.
	flash struct{ Type, Message string }

	// Manifest holds a list of compiled webpack assets.
	manifest map[string]string

	// Data holds template data as a pongo2 context.
	data pongo2.Context
)
