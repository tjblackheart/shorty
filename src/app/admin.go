package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/tjblackheart/shorty/models"
)

func (app App) admin(w http.ResponseWriter, r *http.Request) {
	shorties, err := app.store.Find()
	if err != nil {
		app.err("admin:find", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
	}

	app.render(w, "admin/index.html.j2", data{
		"flash":    app.session.Pop(r.Context(), "flash"),
		"shorties": shorties,
		"_csrf":    csrf.Token(r),
	})
}

func (app App) remove(w http.ResponseWriter, r *http.Request) {
	hashID := mux.Vars(r)["hashID"]
	if err := app.store.DeleteOne(hashID); err != nil {
		app.err("admin:deleteOne", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
	}

	http.Redirect(w, r, "/_a/", 302)
}

func (app App) removeAll(w http.ResponseWriter, r *http.Request) {
	if err := app.store.DeleteMany(); err != nil {
		app.err("admin:removeAll", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
	}

	http.Redirect(w, r, "/_a/", 302)
}

func (app App) importJSON(w http.ResponseWriter, r *http.Request) {
	var shorties []*models.Shorty
	var buf bytes.Buffer

	maxSize := 5 << 20 // 5MB
	actualSize, _ := strconv.Atoi(r.Header.Get("Content-Length"))

	if actualSize > maxSize {
		msg := fmt.Sprintf("Uploaded file too big. The maximum size is %d bytes.", maxSize)
		app.err("admin:import:parse", msg)
		app.session.Put(r.Context(), "flash", flash{"danger", msg})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	// play it safe if the header check didn't work out
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		app.err("admin:import:parseForm", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	file, header, err := r.FormFile("import")
	if err != nil {
		app.err("admin:import:parseForm", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}
	defer file.Close()

	mime := header.Header.Get("Content-Type")
	if mime != "application/json" {
		msg := "Wrong MIME type: " + mime
		app.err("admin:import:mime", msg)
		app.session.Put(r.Context(), "flash", flash{"danger", msg})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	io.Copy(&buf, file)
	if err := json.Unmarshal([]byte(buf.String()), &shorties); err != nil {
		app.err("admin:import:unmarshal", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}
	buf.Reset()

	count, err := app.store.SaveMany(shorties)
	if err != nil {
		app.err("admin:import:saveMany", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	app.session.Put(r.Context(), "flash", flash{"success", fmt.Sprintf("Import successful, %d new entries.", count)})
	http.Redirect(w, r, "/_a/", 302)
}

func (app App) exportJSON(w http.ResponseWriter, r *http.Request) {
	shorties, err := app.store.Find()
	if err != nil {
		app.err("admin:export:find", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
	}

	bs, err := json.Marshal(shorties)
	if err != nil {
		app.err("admin:export:marshal", err.Error())
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	if err := ioutil.WriteFile("/tmp/export.json", bs, 0666); err != nil {
		app.session.Put(r.Context(), "flash", flash{"danger", err.Error()})
		http.Redirect(w, r, "/_a/", 302)
		return
	}

	w.Header().Add("Content-Disposition", "attachment; filename=export.json")
	w.Header().Add("Content-Type", "application/json")

	http.ServeFile(w, r, "/tmp/export.json")
}
