package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flosch/pongo2/v4"
)

func (app App) asset(filename string) string {
	return app.manifest[filename]
}

func (app App) render(w http.ResponseWriter, name string, data pongo2.Context) {
	path := fmt.Sprintf("%s/%s", app.templates, name)
	tpl := pongo2.Must(pongo2.FromFile(path))

	data["asset"] = app.asset

	if err := tpl.ExecuteWriter(data, w); err != nil {
		app.err("render", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) parseManifest() {
	app.manifest = Manifest{}

	bs, err := ioutil.ReadFile("assets/dist/manifest.json")
	if err != nil {
		app.err("parseManifest/readFile", err.Error())
		return
	}

	if err := json.Unmarshal(bs, &app.manifest); err != nil {
		app.err("parseManifest/unmarshal", err.Error())
		return
	}
}
