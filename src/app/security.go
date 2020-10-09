package app

import (
	"net/http"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func (app App) login(w http.ResponseWriter, r *http.Request) {
	flash := app.session.Pop(r.Context(), "flash")
	oldVal := app.session.PopString(r.Context(), "oldVal")

	app.render(w, "page/login.html.j2", Data{
		"flash":  flash,
		"oldVal": oldVal,
		"_csrf":  csrf.Token(r),
	})
}

func (app App) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Destroy(r.Context())
	http.Redirect(w, r, "/", 302)
}

func (app App) authenticate(w http.ResponseWriter, r *http.Request) {
	user := r.PostFormValue("username")
	pass := r.PostFormValue("password")
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
