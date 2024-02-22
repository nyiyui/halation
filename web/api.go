package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	s *Server
}

func OnlyMethodFunc(method string, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, fmt.Sprintf("only %s allowed", method), 405)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) setupAPI() http.Handler {
	r := http.NewServeMux()
	a := &API{s}

	r.Handle("/api/v1/nodes", OnlyMethodFunc("GET", a.nodesGet))
	//r.POST("/node/new", a.nodeNew)
	//r.GET("/node/:node-name", a.nodeGet)
	//r.PATCH("/node/:node-name", a.nodePatch)
	//r.DELETE("/node/:node-name", a.nodeDelete)
	//r.POST("/node/:node-name/activate", a.nodeActivate)
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
