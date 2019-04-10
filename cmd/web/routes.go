package main

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.recoverPanic, app.addHeaders, app.requestLog)

	// add global CSRF
	CSRF := csrf.Protect(
		[]byte(app.config.secret),
		csrf.FieldName("_token"),
		csrf.Secure(!app.config.disableTLS),
	)
	r.Use(CSRF)

	// unsecured
	r.HandleFunc("/", app.index).Methods("GET")
	r.HandleFunc("/", app.submit).Methods("POST")
	r.HandleFunc("/v/{hashID:[a-zA-Z0-9]{3}}", app.view).Methods("GET")
	r.HandleFunc("/{hashID:[a-zA-Z0-9]{3}}", app.redirect).Methods("GET")
	r.HandleFunc("/_l/", app.login).Methods("GET")
	r.HandleFunc("/_l/", app.authenticate).Methods("POST")

	// secured
	a := r.PathPrefix("/_a").Subrouter()
	a.HandleFunc("/", app.admin)
	a.HandleFunc("/r/{hashID:[a-zA-Z0-9]{3}}", app.remove)
	a.HandleFunc("/rm/", app.removeAll).Methods("POST")
	a.HandleFunc("/lg/", app.logout).Methods("POST")
	a.Use(app.securedRoutes)

	// static
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(app.assets)))
	r.PathPrefix("/assets/").Handler(assetsHandler)

	// 404 handler
	r.NotFoundHandler = http.HandlerFunc(app.notFound)

	return r
}
