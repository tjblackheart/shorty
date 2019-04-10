package main

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/tjblackheart/shorty/pkg/models"
)

// 404 error page
func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.gohtml",
		"ui/html/views/404.gohtml",
	}

	app.render(w, r, "base", files, &tplData{})
}

// Form / Index
func (app *application) index(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.gohtml",
		"ui/html/views/index.gohtml",
	}

	app.render(w, r, "base", files, &tplData{})
}

// Calculate shorty
func (app *application) submit(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimSpace(html.EscapeString(r.PostFormValue("url")))

	if validationErr := app.validateInput(url); validationErr != "" {
		session := app.session.Load(r)

		if err := session.PutString(w, "post_data", url); err != nil {
			app.serverError(w, err)
		}

		if err := session.PutString(w, "flash", validationErr); err != nil {
			app.serverError(w, err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// create hashID
	id, err := app.generateHashID(url)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// build model
	s := &models.Shorty{
		Link:    url,
		Shorty:  id,
		Created: time.Now(),
		IP:      r.RemoteAddr,
	}

	// save
	if err := app.shorties.Insert(s); err != nil {
		if err != models.ErrUnique {
			app.serverError(w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/v/%s", s.Shorty), http.StatusFound)
}

// View a shorty
func (app *application) view(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashID := vars["hashID"]
	s, err := app.shorties.Find(hashID)

	if err != nil {
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		}
		app.serverError(w, err)
		return
	}

	files := []string{
		"ui/html/base.gohtml",
		"ui/html/views/shorty.gohtml",
	}

	app.render(w, r, "base", files, &tplData{
		Title:   "View",
		Shorty:  s,
		Request: r,
		Page:    "view",
	})
}

// Redirect to shorty target
func (app *application) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashID := vars["hashID"]
	s, err := app.shorties.Find(hashID)

	if err != nil {
		if err == models.ErrNoRecord {
			app.notFound(w, r)
			return
		}

		app.serverError(w, err)
		return
	}

	if err = app.shorties.AddClick(s); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, s.Link, http.StatusFound)
}

// Remove shorty
func (app *application) remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hashID := vars["hashID"]

	if err := app.shorties.Remove(hashID); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/_a/", http.StatusFound)
}

// Remove all shorties
func (app *application) removeAll(w http.ResponseWriter, r *http.Request) {
	if err := app.shorties.RemoveAll(); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/_a/", http.StatusFound)
}

// Admin list view
func (app *application) admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private")

	shorties, err := app.shorties.FindAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"ui/html/base.gohtml",
		"ui/html/views/admin.gohtml",
	}

	app.render(w, r, "base", files, &tplData{
		Title:    "Admin",
		Shorties: shorties,
		Request:  r,
		Page:     "admin",
	})
}

// Login form
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"ui/html/base.gohtml",
		"ui/html/views/login.gohtml",
	}

	app.render(w, r, "base", files, &tplData{
		Title:   "Login",
		Request: r,
	})
}

// Login form submit
func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	session := app.session.Load(r)
	email := strings.TrimSpace(r.PostFormValue("email"))
	password := strings.TrimSpace(html.EscapeString(r.PostFormValue("password")))
	user, err := app.users.Authenticate(email, password)

	if err != nil {
		if err == models.ErrInvalidCredentials {
			if err = session.PutString(w, "flash", "Invalid credentials"); err != nil {
				app.serverError(w, err)
			}

			if err = session.PutString(w, "last_username", email); err != nil {
				app.serverError(w, err)
			}

			http.Redirect(w, r, "/_l/", http.StatusFound)
			return
		}

		app.serverError(w, err)
		return
	}

	if err = session.PutInt(w, "userID", user.ID); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/_a/", http.StatusFound)
}

// Logout
func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	s := app.session.Load(r)

	if err := s.Destroy(w); err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
