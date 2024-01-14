package web

import (
	"embed"
	"encoding/json"
	"html/template"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").
		Funcs(template.FuncMap{
			"toJSON": json.Marshal,
		}).ParseFS(tmplFS, "templates/*.html"))
}
