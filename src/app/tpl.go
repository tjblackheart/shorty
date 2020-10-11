package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flosch/pongo2/v4"
)

func (app App) initTemplates() {
	app.parseManifest()

	pongo2.DefaultSet.Debug = app.cfg.Debug
	pongo2.DefaultSet.Globals["asset"] = func(filename string) string {
		return app.manifest[filename]
	}
}

func (app *App) parseManifest() {
	app.manifest = manifest{}

	bs, err := ioutil.ReadFile("assets/dist/manifest.json")
	if err != nil {
		app.err("template:parseManifest:readFile", err.Error())
		return
	}

	if err := json.Unmarshal(bs, &app.manifest); err != nil {
		app.err("template:parseManifest:unmarshal", err.Error())
		return
	}
}

func (app App) render(w http.ResponseWriter, name string, d data) {
	path := fmt.Sprintf("%s/%s", app.cfg.ViewsDir, name)
	tpl := pongo2.Must(pongo2.FromCache(path))

	if err := tpl.ExecuteWriter(pongo2.Context(d), w); err != nil {
		app.err("template:render", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
