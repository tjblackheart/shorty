package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tjblackheart/shorty/models"
	"golang.org/x/crypto/bcrypt"
)

func (app App) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, "page/index.html.j2", pongo2.Context{
		"flash": app.session.Pop(r.Context(), "flash"),
		"error": app.session.PopString(r.Context(), "error"),
		"auth":  app.session.GetBool(r.Context(), "_auth"),
		"url":   app.session.PopString(r.Context(), "url"),
		"_csrf": csrf.Token(r),
	})
}

func (app App) create(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(r.PostFormValue("url"))
	s := &models.Shorty{URL: url, IP: r.RemoteAddr}

	if err := s.Validate(); err != nil {
		app.err("page/create/validate", "Validation failed: "+err.Error())
		app.session.Put(r.Context(), "error", "What? This does not look like a valid URL.")
		app.session.Put(r.Context(), "url", url)
		http.Redirect(w, r, "/", 302)
		return
	}

	if exist, _ := app.db.FindOneByURL(s.URL); exist != nil {
		http.Redirect(w, r, fmt.Sprintf("/v/%s", exist.Shorty), 302)
		return
	}

	if err := s.Generate(); err != nil {
		app.err("page/create/generate", "Generator failed: "+err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", "Error generating short link."})
		http.Redirect(w, r, "/", 302)
		return
	}

	if err := app.db.Save(s); err != nil {
		app.err("page/create/add", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", "Error saving data."})
		http.Redirect(w, r, "/", 302)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/v/%s", s.Shorty), 302)
}

func (app App) view(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	shorty, err := app.db.FindOneByShortLink(hashID)

	if err != nil {
		app.err("view/findOne", err.Error())
		app.renderError(w, r, "404", "Link not found. Perhaps it was removed?")
		return
	}

	app.render(w, "page/shorty.html.j2", pongo2.Context{
		"host":      r.Host,
		"shortLink": shorty.Shorty,
	})
}

func (app App) redirect(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	s, err := app.db.FindOneByShortLink(hashID)

	if err != nil {
		app.err("redirect/findOne", err.Error())
		app.renderError(w, r, "404", "Link not found. Perhaps it was removed?")
		return
	}

	s.Clicks++
	if err := app.db.Update(s); err != nil {
		app.err("redirect/update", err.Error())
		app.renderError(w, r, "500", "Something went terribly wrong.")
		return
	}

	http.Redirect(w, r, s.URL, 302)
}

func (app App) login(w http.ResponseWriter, r *http.Request) {
	flash := app.session.Pop(r.Context(), "flash")
	oldVal := app.session.PopString(r.Context(), "oldVal")

	app.render(w, "page/login.html.j2", pongo2.Context{
		"flash":  flash,
		"oldVal": oldVal,
		"_csrf":  csrf.Token(r),
	})
}

func (app App) authenticate(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimSpace(r.PostFormValue("username"))
	pass := strings.TrimSpace(r.PostFormValue("password"))
	app.session.Put(r.Context(), "oldVal", user)

	if user != app.credentials.User {
		app.session.Put(r.Context(), "flash", Flash{"danger", "Invalid credentials."})
		http.Redirect(w, r, "/_l", 302)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(app.credentials.BcryptPass), []byte(pass)); err != nil {
		app.session.Put(r.Context(), "flash", Flash{"danger", "Invalid credentials."})
		http.Redirect(w, r, "/_l", 302)
		return
	}

	app.session.Put(r.Context(), "_auth", true)
	http.Redirect(w, r, "/_a/", 302)
}

func (app App) notFound(w http.ResponseWriter, r *http.Request) {
	app.renderError(w, r, "404", "There is nothing here.")
}

func (app App) renderError(w http.ResponseWriter, r *http.Request, code, msg string) {
	app.render(w, "page/error.html.j2", pongo2.Context{
		"code":    code,
		"message": msg,
	})
}
