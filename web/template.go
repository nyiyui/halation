package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/Masterminds/sprig/v3"
	"nyiyui.ca/halation/node"
)

//go:embed templates
var tmplFS embed.FS
var tmpl *template.Template

//go:embed static
var staticFS embed.FS

func init() {
	tmpl = template.Must(template.New("").
		Funcs(template.FuncMap{
			"partial": func(name string, data interface{}) (interface{}, error) {
				tmpl2, err := tmpl.Clone()
				if err != nil {
					return "", err
				}
				_, err = tmpl2.ParseFS(tmplFS, "templates/"+name)
				if err != nil {
					return "", err
				}
				buf := new(bytes.Buffer)
				err = tmpl2.ExecuteTemplate(buf, name, data)
				if err != nil {
					return "", err
				}
				return template.HTML(buf.String()), nil
			},
			"toJSON": func(v interface{}) ([]byte, error) {
				return json.MarshalIndent(v, "", "  ")
			},
			"nodeTypeFieldlist": func(nodeType string) ([]string, error) {
				// TODO: memoize these for all node types
				newNode, ok := node.NodeTypes[nodeType]
				if !ok {
					return nil, errors.New("invalid node type")
				}
				t := reflect.TypeOf(newNode())
				if t.Kind() == reflect.Pointer {
					t = t.Elem()
				}
				vf := reflect.VisibleFields(t)
				result := make([]string, 0, len(vf))
				for _, field := range vf {
					_, ok := field.Tag.Lookup("halation")
					if !ok {
						continue
					}
					result = append(result, field.Name)
				}
				return result, nil
			},
			"getPromise": func(promises []node.Promise, field string) node.Promise {
				for _, p := range promises {
					if p.FieldName == field {
						return p
					}
				}
				return node.Promise{}
			},
		}).
		Funcs(sprig.FuncMap()).
		ParseFS(tmplFS, "templates/base.html"))
}

func (s *Server) setupStatic() {
	s.sm.Handle("/", http.FileServer(http.FS(staticFS)))
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl2, err := tmpl.Clone()
	if err != nil {
		log.Printf("renderTemplate: %s", err)
		http.Error(w, "rendering template failed", 500)
		return
	}
	_, err = tmpl2.ParseFS(tmplFS, "templates/"+name)
	if err != nil {
		log.Printf("renderTemplate: %s", err)
		http.Error(w, "rendering template failed", 500)
		return
	}
	buf := new(bytes.Buffer)
	err = tmpl2.ExecuteTemplate(buf, name, data)
	if err != nil {
		log.Printf("renderTemplate: %s", err)
		http.Error(w, "rendering template failed", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, buf)
}
