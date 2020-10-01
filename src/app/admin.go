package app

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tjblackheart/shorty/models"
)

func (app App) admin(w http.ResponseWriter, r *http.Request) {
	shorties, err := app.db.Find()
	if err != nil {
		app.err("admin/find", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
	}

	app.render(w, "admin/index.html.j2", pongo2.Context{
		"flash":    app.session.Pop(r.Context(), "flash"),
		"shorties": shorties,
		"_csrf":    csrf.Token(r),
	})
}

func (app App) remove(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	if err := app.db.DeleteOne(hashID); err != nil {
		app.err("admin/deleteOne", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
	}

	http.Redirect(w, r, "/_a/", 302)
}

func (app App) removeAll(w http.ResponseWriter, r *http.Request) {
	if err := app.db.DeleteMany(); err != nil {
		app.err("admin/removeAll", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
	}

	http.Redirect(w, r, "/_a/", 302)
}

func (app App) importJSON(w http.ResponseWriter, r *http.Request) {
	var shorties []*models.Shorty
	var buf bytes.Buffer

	r.ParseMultipartForm(2 << 20) // 2MB
	file, header, err := r.FormFile("import")
	if err != nil {
		app.err("admin/import/parseForm", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}
	defer file.Close()

	mime := header.Header.Get("Content-Type")
	if mime != "application/json" {
		msg := "Wrong MIME type: " + mime
		app.err("admin/import/mime", msg)
		app.session.Put(r.Context(), "flash", Flash{"danger", msg})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	io.Copy(&buf, file)
	if err := json.Unmarshal([]byte(buf.String()), &shorties); err != nil {
		app.err("admin/import/unmarshal", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}
	buf.Reset()

	if err := app.db.SaveMany(shorties); err != nil {
		app.err("admin/import/saveMany", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	app.session.Put(r.Context(), "flash", Flash{"success", "Import ok."})
	http.Redirect(w, r, "/_a/", 302)
}

func (app App) exportJSON(w http.ResponseWriter, r *http.Request) {
	shorties, err := app.db.Find()
	if err != nil {
		app.err("export/find", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
	}

	bs, err := json.Marshal(shorties)
	if err != nil {
		app.err("export/marshal", err.Error())
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	if err := ioutil.WriteFile("/tmp/export.json", bs, 0666); err != nil {
		app.session.Put(r.Context(), "flash", Flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	w.Header().Add("Content-Disposition", "attachment; filename=export.json")
	w.Header().Add("Content-Type", "application/json")

	http.ServeFile(w, r, "/tmp/export.json")
}

func (app App) logout(w http.ResponseWriter, r *http.Request) {
	app.session.Destroy(r.Context())
	http.Redirect(w, r, "/", 302)
}
