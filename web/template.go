package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

//go:embed static
var staticFS embed.FS

func init() {
	tmpl = template.Must(template.New("").
		Funcs(template.FuncMap{
			"toJSON": func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(v, "", "  ")
			},
		}).
		Funcs(sprig.FuncMap()).ParseFS(tmplFS, "templates/*.html"))
}

func (s *Server) setupStatic() {
	s.sm.Handle("/", http.FileServer(http.FS(staticFS)))
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
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, buf)
}
