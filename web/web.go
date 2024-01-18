package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"nyiyui.ca/halation/aiz"
	_ "nyiyui.ca/halation/gradient"
	_ "nyiyui.ca/halation/mpv"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/notify"
	_ "nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/web/tasks"
)

type Change struct {
	NodeName node.NodeName
}

type Server struct {
	sm     *http.ServeMux
	runner *aiz.Runner
	nr     *node.NodeRunner
	tasks  *tasks.Tasks

	changeMuxS *notify.MultiplexerSender[Change]
	changeMux  *notify.Multiplexer[Change]
}

func NewServer(runner *aiz.Runner, nr *node.NodeRunner) *Server {
	s := &Server{
		sm:     http.NewServeMux(),
		runner: runner,
		nr:     nr,
		tasks:  new(tasks.Tasks),
	}
	s.changeMuxS, s.changeMux = notify.NewMultiplexerSender[Change]("Server")
	s.setupStatic()
	s.sm.HandleFunc("/map", s.handleMap)
	s.sm.HandleFunc("/edit", s.handleEdit)
	//s.sm.HandleFunc("/new", s.handleNew)
	//s.sm.HandleFunc("/apply", s.handleApply)
	s.sm.HandleFunc("/tasks", s.handleTasks)
	s.sm.HandleFunc("/events/change", s.handleChange)
	return s
}

func (s *Server) forTemplate(r *http.Request) map[string]interface{} {
	availableNodeTypeNames := make([]string, 0, len(node.NodeTypes))
	for name := range node.NodeTypes {
		availableNodeTypeNames = append(availableNodeTypeNames, name)
	}
	availableStateTypeNames := make([]string, 0, len(aiz.StateTypes))
	for name := range aiz.StateTypes {
		availableStateTypeNames = append(availableStateTypeNames, name)
	}
	availableGradientTypeNames := make([]string, 0, len(aiz.GradientTypes))
	for name := range aiz.GradientTypes {
		availableGradientTypeNames = append(availableGradientTypeNames, name)
	}
	return map[string]interface{}{
		"runner":                     s.runner,
		"nr":                         s.nr,
		"tasks":                      s.tasks,
		"availableNodeTypeNames":     availableNodeTypeNames,
		"availableStateTypeNames":    availableStateTypeNames,
		"availableGradientTypeNames": availableGradientTypeNames,
		"htmx":                       r.Header.Get("HX-Request") == "true",
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.sm.ServeHTTP(w, r)
}

func (s *Server) handleChange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	ch := make(chan Change)
	s.changeMux.Subscribe("handleChange", ch)
	defer s.changeMux.Unsubscribe(ch)
	fmt.Fprint(w, "event: connected\n\n")
	w.(http.Flusher).Flush()
	for cr := range ch {
		fmt.Fprint(w, "event: changed\n")
		data, err := json.Marshal(cr)
		if err != nil {
			log.Printf("handleChange: marshal: %s", err)
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		w.(http.Flusher).Flush()
	}
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET allowed", 405)
		return
	}
	s.renderTemplate(w, r, "tasks.html", map[string]interface{}{
		"s": s.forTemplate(r),
	})
}

func (s *Server) handleMap(w http.ResponseWriter, r *http.Request) {
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	listeners := s.nr.NM.GenListeners()
	roots := make([]node.NodeName, 0)
	for key, node := range s.nr.NM.Nodes {
		if len(node.GetListensTo()) == 0 {
			roots = append(roots, key)
		}
	}
	data := s.forTemplate(r)
	data["roots"] = roots
	data["listeners"] = listeners
	s.renderTemplate(w, r, "nodemap.html", data)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "index.html", map[string]interface{}{
		"runner": s.runner,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	io.Copy(w, buf)
}
