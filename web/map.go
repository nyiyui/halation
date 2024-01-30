package web

import (
	"net/http"

	"nyiyui.ca/halation/node"
)

func (s *Server) handleMap2(w http.ResponseWriter, r *http.Request) {
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	pm := s.nr.NM.GenPromiseMap()
	roots := make([]node.NodeName, 0)
	for key, node := range s.nr.NM.Nodes {
		if len(node.BaseNodeRef().Promises) == 0 {
			roots = append(roots, key)
		}
	}
	data := s.forTemplate(r)
	data["roots"] = roots
	data["pm"] = pm
	s.renderTemplate(w, r, "nodemap.html", data)
}
