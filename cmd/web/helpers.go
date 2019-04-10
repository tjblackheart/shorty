package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/speps/go-hashids"
)

// error logger
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.err.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// client error
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// check auth session
func (app *application) isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	sess := app.session.Load(r)
	userID, err := sess.GetInt("userID")

	if err != nil {
		app.serverError(w, err)
	}

	if userID == 0 {
		return false
	}

	return true
}

// validate input
func (app *application) validateInput(input string) string {
	e := ""
	u, err := url.Parse(input)

	if err != nil || u.Scheme == "" || u.Host == "" || u.Scheme != "http" && u.Scheme != "https" {
		e = "What? This has to be an URL, you know."
	}

	return e
}

func (app *application) generateHashID(url string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = fmt.Sprintf("%f", rand.ExpFloat64()*1e9)
	hd.MinLength = 3

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	id, err := h.Encode([]int{len(url)})
	if err != nil {
		return "", err
	}

	return id, nil
}
