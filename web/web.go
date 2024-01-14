package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/mpv"
)

type Server struct {
	sm *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		sm: http.NewServeMux(),
	}
	s.sm.HandleFunc("/", s.handleIndex)
	s.sm.HandleFunc("/edit", s.handleEdit)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.sm.ServeHTTP(w, r)
}

func (s *Server) handleEdit(w http.ResponseWriter, r *http.Request) {
	runner := &aiz.Runner{Specific: map[string]interface{}{}}
	show := &aiz.Show{Cues: []aiz.Cue{
		{Name: "0 paused", SGs: []aiz.SG{
			{State: &mpv.State{
				FilePath:   "./big_buck_bunny_480p_h264.mov",
				Paused:     mpv.Ptr(true),
				Position:   mpv.Ptr(0),
				Fullscreen: mpv.Ptr(false),
			}},
		}},
		{Name: "1 playing", SGs: []aiz.SG{
			{State: &mpv.State{
				FilePath:   "./big_buck_bunny_480p_h264.mov",
				Paused:     mpv.Ptr(false),
				Position:   mpv.Ptr(60),
				Fullscreen: mpv.Ptr(false),
			}},
		}},
	}}

	cueIndex, err := strconv.ParseInt(r.URL.Query().Get("cue-i"), 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("parse cue-i: %s", err), 422)
		return
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "edit.html", map[string]interface{}{
		"runner":   runner,
		"show":     show,
		"cueIndex": cueIndex,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	io.Copy(w, buf)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	runner := &aiz.Runner{Specific: map[string]interface{}{}}
	show := &aiz.Show{Cues: []aiz.Cue{
		{Name: "0 paused", SGs: []aiz.SG{
			{State: &mpv.State{
				FilePath:   "./big_buck_bunny_480p_h264.mov",
				Paused:     mpv.Ptr(true),
				Position:   mpv.Ptr(0),
				Fullscreen: mpv.Ptr(false),
			}},
		}},
		{Name: "1 playing", SGs: []aiz.SG{
			{State: &mpv.State{
				FilePath:   "./big_buck_bunny_480p_h264.mov",
				Paused:     mpv.Ptr(false),
				Position:   mpv.Ptr(60),
				Fullscreen: mpv.Ptr(false),
			}},
		}},
	}}
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "index.html", map[string]interface{}{
		"runner": runner,
		"show":   show,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), 500)
		return
	}
	io.Copy(w, buf)
}
