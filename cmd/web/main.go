package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tjblackheart/shorty/pkg/models/sqlite"
)

type application struct {
	config struct {
		host, dsn, secret string
		cert, key         string
		disableTLS        bool
	}
	assets    string
	info, err *log.Logger
	shorties  *sqlite.ShortyModel
	users     *sqlite.UserModel
	session   *scs.Manager
}

var app application
var tlsConfig = &tls.Config{
	PreferServerCipherSuites: true,
	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	CipherSuites: []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	},
	MinVersion: tls.VersionTLS12,
}

func init() {
	app.info = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.err = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.assets = "./ui/static/"

	flag.StringVar(&app.config.dsn, "dsn", "./var/shorty.db", "Database DSN")
	flag.BoolVar(&app.config.disableTLS, "disableTLS", false, "Start without TLS")
	flag.StringVar(&app.config.cert, "cert", "./tls/cert.pem", "Path to cert.pem")
	flag.StringVar(&app.config.key, "key", "./tls/key.pem", "Path to key.pem")
	flag.Parse()

	app.config.host = os.Getenv("APP_HOST")
	if app.config.host == "" {
		app.info.Println("Missing APP_HOST, using default value of :3000")
		app.config.host = ":3000"
	}

	app.config.secret = os.Getenv("APP_SECRET")
	if app.config.secret == "" {
		app.err.Fatalln("Missing APP_SECRET")
	}

	app.session = scs.NewCookieManager(app.config.secret)
	app.session.Lifetime(time.Hour)
	app.session.SameSite("Strict")

	if !app.config.disableTLS {
		app.session.Secure(true) // https
	}
}

func main() {
	db, err := openDB(app.config.dsn)
	if err != nil {
		app.err.Fatal(err)
	}
	defer db.Close()

	// init tables after DB init
	app.shorties = &sqlite.ShortyModel{DB: db}
	if err = app.shorties.InitTable(); err != nil {
		app.err.Fatal(err)
	}

	app.users = &sqlite.UserModel{DB: db}
	if err = app.users.InitTable(); err != nil {
		app.err.Fatal(err)
	}

	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         app.config.host,
		ErrorLog:     app.err,
		IdleTimeout:  time.Minute,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		TLSConfig:    tlsConfig,
	}

	app.info.Printf("Starting server on %s\n", app.config.host)
	app.info.Printf("To create a user use bin/create_user.\n")
	app.info.Printf("use bin/shorty -h to see all startup options.\n")

	if app.config.disableTLS {
		app.err.Fatal(srv.ListenAndServe())
	} else {
		app.err.Fatal(srv.ListenAndServeTLS(app.config.cert, app.config.key))
	}
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
