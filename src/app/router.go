package app

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func (app App) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.recover, app.requestLog, app.loadSession)

	r.Use(csrf.Protect(
		[]byte(app.cfg.Secret),
		csrf.FieldName("_csrf"),
		csrf.Secure(false),
		csrf.SameSite(csrf.SameSiteStrictMode),
	))

	r.HandleFunc("/", app.index).Methods(http.MethodGet)
	r.HandleFunc("/", app.create).Methods(http.MethodPost)
	r.HandleFunc("/_l", app.login).Methods(http.MethodGet)
	r.HandleFunc("/_l", app.authenticate).Methods(http.MethodPost)
	r.HandleFunc("/{hashID:[a-zA-Z0-9]{3}}", app.redirect)
	r.HandleFunc("/v/{hashID:[a-zA-Z0-9]{3}}", app.view)

	r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/dist"))))
	r.NotFoundHandler = http.HandlerFunc(app.notFound)

	a := r.PathPrefix("/_a").Subrouter()
	a.Use(app.authenticated)

	a.HandleFunc("/", app.admin).Methods(http.MethodGet)
	a.HandleFunc("/r/{hashID:[a-zA-Z0-9]{3}}", app.removeSingle).Methods(http.MethodGet)
	a.HandleFunc("/r", app.removeAll).Methods(http.MethodGet)
	a.HandleFunc("/i", app.importJSON).Methods(http.MethodPost)
	a.HandleFunc("/e", app.exportJSON).Methods(http.MethodGet)
	a.HandleFunc("/l", app.logout).Methods(http.MethodGet)

	return r
}
