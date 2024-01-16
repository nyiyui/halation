package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/web/tasks"
)

type Server struct {
	sm     *http.ServeMux
	runner *aiz.Runner
	nr     *node.NodeRunner
	tasks  *tasks.Tasks
}

func NewServer(runner *aiz.Runner, show *aiz.Show) *Server {
	s := &Server{
		sm:     http.NewServeMux(),
		runner: runner,
		show:   show,
		tasks:  new(tasks.Tasks),
	}
	s.sm.HandleFunc("/", s.handleIndex)
	s.sm.HandleFunc("/edit", s.handleEdit)
	//s.sm.HandleFunc("/apply", s.handleApply)
	s.sm.HandleFunc("/tasks", s.handleTasks)
	//s.sm.HandleFunc("/cue-applied", s.handleCueApplied)
	return s
}

func (s *Server) forTemplate() map[string]interface{} {
	return map[string]interface{}{
		"runner": s.runner,
		"show":   s.show,
		"tasks":  s.tasks,
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

//func (s *Server) handleApply(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		http.Error(w, "only POST allowed", 405)
//		return
//	}
//	cueIndex, err := strconv.ParseInt(r.URL.Query().Get("cue-i"), 10, 32)
//	if err != nil {
//		http.Error(w, fmt.Sprintf("parse cue-i: %s", err), 422)
//		return
//	}
//	s.tasks.Append(fmt.Sprintf("apply cue %d", cueIndex), func() error {
//		return s.show.ApplyCue(s.runner, int(cueIndex), context.Background())
//	})
//	s.renderTemplate(w, r, "apply.html", map[string]interface{}{
//		"runner":   s.runner,
//		"show":     s.show,
//		"cueIndex": cueIndex,
//	})
//}

func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	cueIndex, err := strconv.ParseInt(r.URL.Query().Get("cue-i"), 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse cue-i: %s", err), 422)
		return
	}
	s.renderTemplate(w, r, "edit.html", map[string]interface{}{
		"runner":   s.runner,
		"show":     s.show,
		"cueIndex": cueIndex,
	})
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "index.html", map[string]interface{}{
		"runner": s.runner,
		"show":   s.show,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	io.Copy(w, buf)
}
