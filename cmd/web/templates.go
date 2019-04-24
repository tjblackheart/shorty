package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/tjblackheart/shorty/pkg/models"
)

type tplData struct {
	Title              string
	Request            *http.Request
	Shorty             *models.Shorty
	Shorties           []*models.Shorty
	PostData, Flash    string
	LastUsername, Page string
	UserID             int
	CSRFToken          string
}

// template helper functions map
var functions = template.FuncMap{
	"formatDate": formatDate,
	"formatLink": formatLink,
	"shorten":    shorten,
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return "?"
	}

	return t.Format("Jan 02 2006, 15:04")
}

func formatLink(s string, r *http.Request) string {
	scheme := "https"

	if r.TLS == nil {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s/%s", scheme, r.Host, s)
}

func shorten(s string, n int) string {
	if len(s) > n {
		return s[:n] + " ... "
	}

	return s
}

// the render function
func (app *application) render(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	files []string,
	data *tplData,
) {
	buf := new(bytes.Buffer)

	// register functions
	ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
	}

	// add data
	session := app.session.Load(r)
	data.Flash, _ = session.PopString(w, "flash")
	data.LastUsername, _ = session.PopString(w, "last_username")
	data.PostData, _ = session.PopString(w, "post_data")
	data.UserID, _ = session.GetInt("userID")
	data.CSRFToken = csrf.Token(r)

	// render to buffer first to catch render errors in partials
	if err = ts.ExecuteTemplate(buf, name, &data); err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}
