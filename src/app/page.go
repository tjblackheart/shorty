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
		"_csrf": csrf.Token(r),
	})
}

func (app App) create(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	s := &models.Shorty{URL: url, IP: r.RemoteAddr}

	if err := s.Validate(); err != nil {
		app.err("page/create/validate", "Validation failed: "+err.Error())
		app.session.Put(r.Context(), "error", "What? This does not look like a valid URL.")
		http.Redirect(w, r, "/", 302)
		return
	}

	if exist, _ := app.repo.FindOneByURL(s.URL); exist != nil {
		http.Redirect(w, r, fmt.Sprintf("/v/%s", exist.Shorty), 302)
		return
	}

	if err := s.Generate(); err != nil {
		app.err("page/create/generate", "Generator failed: "+err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", "Error generating short link."})
		http.Redirect(w, r, "/", 302)
		return
	}

	if err := app.repo.Save(s); err != nil {
		app.err("page/create/add", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", "Error saving data."})
		http.Redirect(w, r, "/", 302)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/v/%s", s.Shorty), 302)
}

func (app App) view(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	shorty, err := app.repo.FindOneByShortLink(hashID)

	if err != nil {
		app.err("view/findOne", err.Error())
		app.session.Put(r.Context(), "message", "Link not found. Perhaps it was removed?")
		http.Redirect(w, r, "/error/404", 302)
		return
	}

	app.render(w, "page/shorty.html.j2", pongo2.Context{
		"host":      r.Host,
		"shortLink": shorty.Shorty,
	})
}

func (app App) redirect(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	s, err := app.repo.FindOneByShortLink(hashID)

	if err != nil {
		app.err("redirect/findOne", err.Error())
		app.session.Put(r.Context(), "message", "Link not found. Perhaps it was removed?")
		http.Redirect(w, r, "/error/404", 302)
		return
	}

	s.Clicks++
	if err := app.repo.Update(s); err != nil {
		app.err("redirect/update", err.Error())
		http.Redirect(w, r, "/error/500", 500)
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

func (app App) renderError(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	msg := app.session.PopString(r.Context(), "message")

	if msg == "" {
		msg = "Something went terribly wrong."
	}

	app.render(w, "page/error.html.j2", pongo2.Context{
		"code":    code,
		"message": msg,
	})
}

func (app App) render404(w http.ResponseWriter, r *http.Request) {
	app.render(w, "page/404.html.j2", pongo2.Context{})
}
