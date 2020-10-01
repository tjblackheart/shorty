package app

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (app App) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s %s", r.RemoteAddr, r.UserAgent(), r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app App) loadSession(next http.Handler) http.Handler {
	return app.session.LoadAndSave(next)
}

func (app App) authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.session.GetBool(r.Context(), "_auth") {
			http.Redirect(w, r, "/_l", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app App) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.err("recover", fmt.Sprintf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
