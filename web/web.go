package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"nyiyui.ca/halation/aiz"
	_ "nyiyui.ca/halation/gradient"
	_ "nyiyui.ca/halation/mpv"
	"nyiyui.ca/halation/node"
	_ "nyiyui.ca/halation/osc"
	"nyiyui.ca/halation/web/tasks"
)

type Server struct {
	sm     *http.ServeMux
	runner *aiz.Runner
	nr     *node.NodeRunner
	tasks  *tasks.Tasks
}

func NewServer(runner *aiz.Runner, nr *node.NodeRunner) *Server {
	s := &Server{
		sm:     http.NewServeMux(),
		runner: runner,
		nr:     nr,
		tasks:  new(tasks.Tasks),
	}
	s.setupStatic()
	s.sm.HandleFunc("/map", s.handleMap)
	s.sm.HandleFunc("/edit", s.handleEdit)
	//s.sm.HandleFunc("/new", s.handleNew)
	//s.sm.HandleFunc("/apply", s.handleApply)
	s.sm.HandleFunc("/tasks", s.handleTasks)
	//s.sm.HandleFunc("/cue-applied", s.handleCueApplied)
	return s
}

func (s *Server) forTemplate() map[string]interface{} {
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
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.sm.ServeHTTP(w, r)
}

//func (s *Server) handleCueApplied(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Cache-Control", "no-cache")
//	w.Header().Set("Connection", "keep-alive")
//	w.Header().Set("Content-Type", "text/event-stream")
//	ch := make(chan aiz.CueRequest)
//	s.runner.CueAppliedMux.Subscribe("handleCueApplied", ch)
//	defer s.runner.CueAppliedMux.Unsubscribe(ch)
//	fmt.Fprint(w, "event: connected\n\n")
//	w.(http.Flusher).Flush()
//	for cr := range ch {
//		fmt.Fprint(w, "event: cue-applied\n")
//		data, err := json.Marshal(cr)
//		if err != nil {
//			log.Printf("handleCueApplied: marshal: %s", err)
//			continue
//		}
//		fmt.Fprintf(w, "data: %s\n\n", data)
//		w.(http.Flusher).Flush()
//	}
//}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET allowed", 405)
		return
	}
	s.renderTemplate(w, r, "tasks.html", map[string]interface{}{
		"s": s.forTemplate(),
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
	data := s.forTemplate()
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
