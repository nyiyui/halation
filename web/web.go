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
	s.sm.HandleFunc("/", s.handleIndex)
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

func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	if !(r.Method == "GET" || r.Method == "POST") {
		http.Error(w, "only GET, POST allowed", 405)
		return
	}

	nodeName := node.ParseNodeName(r.URL.Query().Get("node-name"))
	ok := func() bool {
		s.nr.NMLock.RLock()
		defer s.nr.NMLock.RUnlock()
		_, ok := s.nr.NM.Nodes[nodeName]
		if !ok {
			http.Error(w, fmt.Sprintf("node %s does not exist", nodeName), 404)
			return false
		}
		return true
	}()
	if !ok {
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("parsing form data failed: %s", err)
			http.Error(w, "parsing form data failed", 422)
			return
		}
		newNodeFn, ok := node.NodeTypes[r.PostForm.Get("type")]
		if !ok {
			http.Error(w, "invalid node type", 422)
			return
		}
		node2 := newNodeFn()
		switch node2.(type) {
		case *node.Manual:
		case *node.SetState:
			// === State
			newStateFn, ok := aiz.StateTypes[r.PostForm.Get("state-type")]
			if !ok {
				fmt.Println(r.PostForm.Get("state-type"))
				fmt.Println(aiz.StateTypes)
				http.Error(w, "invalid state type", 422)
				return
			}
			state2 := newStateFn()
			err := json.Unmarshal([]byte(r.PostForm.Get("state")), state2)
			if !ok {
				http.Error(w, fmt.Sprintf("parsing state JSON failed: %s", err), 422)
				return
			}

			// === Gradient
			gradientType := r.PostForm.Get("gradient-type")
			var gradient2 aiz.Gradient
			if gradientType != "" {
				newGradientFn, ok := aiz.GradientTypes[r.PostForm.Get("gradient-type")]
				if !ok {
					fmt.Println(r.PostForm.Get("gradient-type"))
					fmt.Println(aiz.GradientTypes)
					http.Error(w, "invalid gradient type", 422)
					return
				}
				gradient2 = newGradientFn()
				err = json.Unmarshal([]byte(r.PostForm.Get("gradient")), gradient2)
				if !ok {
					http.Error(w, fmt.Sprintf("parsing gradient JSON failed: %s", err), 422)
					return
				}
			}

			node2.(*node.SetState).SG = &aiz.SG{
				State:    state2,
				Gradient: gradient2,
			}
		default:
			http.Error(w, "node type not implemented", 422)
			return
		}
		listensTo := make([]node.NodeName, 0)
		err = json.Unmarshal([]byte(r.PostForm.Get("listens-to")), &listensTo)
		if !ok {
			http.Error(w, fmt.Sprintf("parsing listens-to JSON failed: %s", err), 422)
			return
		}
		node2.SetListensTo(listensTo)
		func() {
			s.nr.NMLock.Lock()
			defer s.nr.NMLock.Unlock()
			s.nr.NM.Nodes[nodeName] = node2
		}()
	}
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	data := s.forTemplate()
	data["node"] = s.nr.NM.Nodes[nodeName]
	data["nodeName"] = nodeName
	data["saved"] = r.Method == "POST"
	s.renderTemplate(w, r, "edit.html", data)
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
