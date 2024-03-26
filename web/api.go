package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"nyiyui.ca/halation/node"
)

type API struct {
	s *Server
}

func (s *Server) setupAPI() http.Handler {
	r := http.NewServeMux()
	a := &API{s}

	r.HandleFunc("GET /api/v1/nodes", a.nodesGet)
	r.HandleFunc("POST /api/v1/node/{nodeName}", a.nodeNew)
	//r.GET("/node/:node-name", a.nodeGet)
	r.HandleFunc("PATCH /api/v1/node/{nodeName}", a.nodePatch)
	//r.DELETE("/node/:node-name", a.nodeDelete)
	r.HandleFunc("POST /api/v1/node/{nodeName}/activate", a.nodeActivate)
	//r.GET("/nodes/events", a.nodesEvents)

	//r.POST("/cues", a.cuesGet)
	//r.GET("/cues", a.cuesGet)
	//r.GET("/cues/events", a.cuesEvents)

	//r.GET("/export", a.handleExport)
	//r.POST("/import", a.handleImport)
	return r
}

func (a *API) nodesGet(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(a.s.nr.NM)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

func (a *API) nodeNew(w http.ResponseWriter, r *http.Request) {
	var nj node.NodeJSON
	err := json.NewDecoder(r.Body).Decode(&nj)
	if err != nil {
		http.Error(w, fmt.Sprintf("decode json body: %s", err), 422)
		return
	}
	nodeName := node.ParseNodeName(r.PathValue("nodeName"))
	ok := func() bool {
		a.s.nr.NMLock.RLock()
		defer a.s.nr.NMLock.RUnlock()
		_, ok := a.s.nr.NM.Nodes[nodeName]
		if ok {
			http.Error(w, fmt.Sprintf("node %s already exists", nodeName), 404)
			return false
		}
		a.s.nr.NM.Nodes[nodeName] = nj.Node
		return true
	}()
	if !ok {
		return
	}
	a.s.autosave()
	http.Error(w, "{}", 200)
}

func (a *API) nodePatch(w http.ResponseWriter, r *http.Request) {
	var nj node.NodeJSON
	err := json.NewDecoder(r.Body).Decode(&nj)
	if err != nil {
		http.Error(w, fmt.Sprintf("decode json body: %s", err), 422)
		return
	}
	nodeName := node.ParseNodeName(r.PathValue("nodeName"))
	ok := func() bool {
		a.s.nr.NMLock.RLock()
		defer a.s.nr.NMLock.RUnlock()
		_, ok := a.s.nr.NM.Nodes[nodeName]
		if !ok {
			http.Error(w, fmt.Sprintf("node %s does not exist", nodeName), 404)
			return false
		}
		a.s.nr.NM.Nodes[nodeName] = nj.Node
		return true
	}()
	if !ok {
		return
	}
	a.s.autosave()
	http.Error(w, "{}", 200)
}

func (a *API) nodeActivate(w http.ResponseWriter, r *http.Request) {
	nodeName := node.ParseNodeName(r.PathValue("nodeName"))
	ok := func() bool {
		a.s.nr.NMLock.RLock()
		defer a.s.nr.NMLock.RUnlock()
		_, ok := a.s.nr.NM.Nodes[nodeName]
		if !ok {
			http.Error(w, fmt.Sprintf("node %s does not exist", nodeName), 404)
			return false
		}
		return true
	}()
	if !ok {
		return
	}
	a.s.nr.ActivateNodeUsingPromises(nodeName, nil)
	http.Error(w, "{}", 200)
}
