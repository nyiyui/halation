package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/gradient"
	"nyiyui.ca/halation/node"
	"nyiyui.ca/halation/timeutil"
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

		if typeNew := r.PostForm.Get("type"); typeNew != r.PostForm.Get("type-original") {
			// type changed, redirect to page with node-type-override
			u := *r.URL
			q := u.Query()
			q.Set("node-type-override", typeNew)
			u.RawQuery = q.Encode()
			http.Redirect(w, r, u.String(), 302)
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
		case *node.EvalLua:
			source := r.PostForm.Get("source")
			node2.(*node.EvalLua).Source = source
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
				switch gradientType {
				case "nyiyui.ca/halation/gradient.LinearGradient":
					gradient2 := gradient2.(*gradient.LinearGradient)
					duration, err := time.ParseDuration(r.PostForm.Get("gradient-duration"))
					if err != nil {
						http.Error(w, fmt.Sprintf("parsing gradient form failed: duration: %s", err), 422)
						return
					}
					gradient2.Duration_ = timeutil.Duration(duration)
					preferredResolution, err := time.ParseDuration(r.PostForm.Get("gradient-preferred-resolution"))
					if err != nil {
						http.Error(w, fmt.Sprintf("parsing gradient form failed: preferred-resolution: %s", err), 422)
						return
					}
					gradient2.PreferredResolution_ = timeutil.Duration(preferredResolution)
				default:
					err = json.Unmarshal([]byte(r.PostForm.Get("gradient")), gradient2)
					if !ok {
						http.Error(w, fmt.Sprintf("parsing gradient JSON failed: %s", err), 422)
						return
					}
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
		promises, err := parseFormPromises(r.PostForm)
		if err != nil {
			http.Error(w, fmt.Sprintf("parsing promises failed: %s", err), 422)
			return
		}
		node2.BaseNodeRef().Promises = promises
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

	n := s.nr.NM.Nodes[nodeName]
	nodeTypeOverride := r.URL.Query().Get("node-type-override")
	if nodeTypeOverride != "" {
		newNode, ok := node.NodeTypes[nodeTypeOverride]
		if !ok {
			http.Error(w, "node-type-override: invalid node type", 422)
			return
		}
		n2 := newNode()
		*n2.BaseNodeRef() = *n.BaseNodeRef()
		n = n2
	}
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	data := s.forTemplate(r)
	data["node"] = n
	data["nodeName"] = nodeName
	s.renderTemplate(w, r, "edit.html", data)
}

func parseFormPromises(data url.Values) ([]node.Promise, error) {
	p := map[string]node.Promise{}
	enabled := []string{}
	for name, value := range data {
		if !strings.HasPrefix(name, "promise-") {
			continue
		}
		parts := strings.Split(name, "-")
		if len(parts) != 3 {
			return nil, fmt.Errorf("%s: must have 3 parts only", name)
		}
		id := parts[1]
		switch parts[2] {
		case "enable":
			enabled = append(enabled, id)
		case "field":
			n := p[id]
			n.FieldName = value[0]
			p[id] = n
		case "supply":
			n := p[id]
			n.SupplyNodeName = node.ParseNodeName(value[0])
			p[id] = n
		}
	}
	result := make([]node.Promise, 0, len(enabled))
	for _, id := range enabled {
		result = append(result, p[id])
	}
	return result, nil
}
