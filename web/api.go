package web

import "github.com/julienschmidt/httprouter"

type API struct {
	s *Server
}

func (s *Server) setupAPI() http.Handler {
	r := httprouter.New()
	a := &API{s}

	r.GET("/nodes", a.nodesGet)
	r.POST("/node/new", a.nodeNew)
	r.GET("/node/:node-name", a.nodeGet)
	r.PATCH("/node/:node-name", a.nodePatch)
	r.DELETE("/node/:node-name", a.nodeDelete)
	r.POST("/node/:node-name/activate", a.nodeActivate)
	r.GET("/nodes/events", a.nodesEvents)

	r.POST("/cues", a.cuesGet)
	r.GET("/cues", a.cuesGet)
	r.GET("/cues/events", a.cuesEvents)

	r.GET("/export", a.handleExport)
	r.POST("/import", a.handleImport)
	return sm
}
