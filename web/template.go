package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
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

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	buf := new(bytes.Buffer)
	var err error
	err = tmpl.ExecuteTemplate(buf, name, data)
	if err != nil {
		log.Printf("renderTemplate: %s", err)
		http.Error(w, "rendering template failed", 500)
		return
	}
	io.Copy(w, buf)
}
