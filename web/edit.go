package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/node"
)

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

		origNodeName := nodeName
		nodeName = node.ParseNodeName(r.PostForm.Get("name"))

		newNodeFn, ok := node.NodeTypes[r.PostForm.Get("type")]
		if !ok {
			http.Error(w, fmt.Sprintf("invalid node type %s", r.PostForm.Get("type")), 422)
			return
		}
		node2 := newNodeFn()
		node2.SetDescription(r.PostForm.Get("description"))
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
			s.changeMuxS.Send(Change{NodeName: nodeName})
			s.nr.NMLock.Lock()
			defer s.nr.NMLock.Unlock()
			s.nr.NM.Nodes[nodeName] = node2
			if origNodeName != nodeName {
				delete(s.nr.NM.Nodes, origNodeName)
			}

			u := &url.URL{Path: "/edit"}
			v := url.Values{}
			v.Add("node-name", nodeName.String())
			u.RawQuery = v.Encode()
			http.Redirect(w, r, u.String(), 302)
			return
		}()
	}
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	data := s.forTemplate(r)
	data["node"] = s.nr.NM.Nodes[nodeName]
	data["nodeName"] = nodeName
	data["htmx"] = true
	s.renderTemplate(w, r, "edit.html", data)
}
