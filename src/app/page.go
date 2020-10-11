package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tjblackheart/shorty/models"
)

func (app App) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, "page/index.html.j2", data{
		"flash": app.session.Pop(r.Context(), "flash"),
		"error": app.session.PopString(r.Context(), "error"),
		"auth":  app.session.GetBool(r.Context(), "_auth"),
		"url":   app.session.PopString(r.Context(), "url"),
		"_csrf": csrf.Token(r),
	})
}

func (app App) create(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}

	url := strings.TrimSpace(r.PostFormValue("url"))
	url = app.policy.Sanitize(url)

	s := &models.Shorty{URL: url, IP: ip}
	if err := s.Validate(); err != nil {
		app.err("page:create:validate", fmt.Sprintf("Validation failed: %s: %s", err.Error(), s.URL))
		app.session.Put(r.Context(), "error", "What? This does not look like a valid URL.")
		app.session.Put(r.Context(), "url", url)
		http.Redirect(w, r, "/", 302)
		return
	}

	if exist, _ := app.store.FindOneByURL(s.URL); exist != nil {
		http.Redirect(w, r, fmt.Sprintf("/v/%s", exist.Shorty), 302)
		return
	}

	if err := s.Generate(); err != nil {
		app.err("page:create:generate", "Generator failed: "+err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", "Error generating short link."})
		http.Redirect(w, r, "/", 302)
		return
	}

	if err := app.store.Save(s); err != nil {
		app.err("page:create:add", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", "Error saving data."})
		http.Redirect(w, r, "/", 302)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/v/%s", s.Shorty), 302)
}

func (app App) view(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	shorty, err := app.store.FindOneByShortLink(hashID)

	if err != nil {
		app.err("page:view:findOne", err.Error())
		app.renderError(w, r, "404", "Link not found. Perhaps it was removed?")
		return
	}

	app.render(w, "page/shorty.html.j2", data{
		"host":      r.Host,
		"shortLink": shorty.Shorty,
	})
}

func (app App) redirect(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	s, err := app.store.FindOneByShortLink(hashID)

	if err != nil {
		app.err("page:redirect:findOne", err.Error())
		app.renderError(w, r, "404", "Link not found. Perhaps it was removed?")
		return
	}

	s.Clicks++
	if err := app.store.Update(s); err != nil {
		app.err("page:redirect:update", err.Error())
		app.renderError(w, r, "500", "Something went terribly wrong.")
		return
	}

	http.Redirect(w, r, s.URL, 302)
}

func (app App) notFound(w http.ResponseWriter, r *http.Request) {
	app.renderError(w, r, "404", "There is nothing here.")
}

func (app App) renderError(w http.ResponseWriter, r *http.Request, code, msg string) {
	app.render(w, "page/error.html.j2", data{
		"code":    code,
		"message": msg,
	})
}
