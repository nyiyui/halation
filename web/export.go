package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"nyiyui.ca/halation/node"
)

type Export struct {
	Metadata ExportMetadata
	Cuelist  *node.Cuelist
	NodeMap  *node.NodeMap
}

type ExportMetadata struct {
	ExportGeneratedAt time.Time
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	s.nr.NMLock.RLock()
	defer s.nr.NMLock.RUnlock()
	e := Export{
		Metadata: ExportMetadata{
			ExportGeneratedAt: time.Now(),
		},
		Cuelist: s.cuelist,
		NodeMap: s.nr.NM,
	}
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("export: marshal failed: %s", err)
		http.Error(w, "marshal failed", 500)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
