package mpv

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/mpv"
	"nyiyui.ca/halation/aiz"
)

type State struct {
	FilePath string
	Paused   *bool
	// TODO: Position doesn't work when the file is initially loaded.
	Position        *int
	Fullscreen      *bool
	ExtraProperties map[string]interface{} `json:"extraProperties"`
}

func (s *State) Clone() aiz.State {
	s2 := State{
		FilePath: s.FilePath,
	}
	if s.Paused != nil {
		paused := *s.Paused
		s2.Paused = &paused
	}
	if s.Position != nil {
		position := *s.Position
		s2.Position = &position
	}
	if s.Fullscreen != nil {
		fullscreen := *s.Fullscreen
		s2.Fullscreen = &fullscreen
	}
	s2.ExtraProperties = map[string]interface{}{}
	for key, val := range s.ExtraProperties {
		s2.ExtraProperties[key] = val
	}
	return &s2
}

func (s *State) String() string {
	b := new(strings.Builder)
	fmt.Fprintf(b, "play %s", s.FilePath)
	if s.Paused != nil {
		if *s.Paused {
			b.WriteString(" (pause)")
		} else {
			b.WriteString(" (play)")
		}
	}
	if s.Position != nil {
		fmt.Fprintf(b, " (%ds)", *s.Position)
	}
	if s.Fullscreen != nil {
		if *s.Fullscreen {
			b.WriteString(" (full)")
		} else {
			b.WriteString(" (window)")
		}
	}
	for key, value := range s.ExtraProperties {
		fmt.Fprintf(b, " (%s=%s)", key, value)
	}
	return b.String()
}

func (s *State) Reify(r *aiz.Runner, g aiz.Gradient, prev_ aiz.State) error {
	c_, ok := r.Specific[packageName]
	if !ok {
		return errors.New("runner doesn't have mpv client")
	}
	c, ok := c_.(*Client)
	if !ok {
		return errors.New("runner has wrong client")
	}
	_, ok = prev_.(*State)
	if !ok && prev_ != nil {
		return errors.New("prev isn't of same type")
	}

	// filepath
	currentPath, err := c.client.Path()
	if err != nil {
		return err
	}
	if s.FilePath != "" && currentPath != s.FilePath {
		err = c.client.Loadfile(s.FilePath, mpv.LoadFileModeReplace)
		if err != nil {
			return err
		}
	}

	if s.Paused != nil {
		err = c.client.SetPause(*s.Paused)
		if err != nil {
			return err
		}
	}

	if s.Position != nil {
		err = c.client.Seek(*s.Position, mpv.SeekModeAbsolute)
		if err != nil {
			return err
		}
	}

	if s.Fullscreen != nil {
		err = c.client.SetFullscreen(*s.Fullscreen)
		if err != nil {
			return err
		}
	}

	for name, value := range s.ExtraProperties {
		err = c.client.SetProperty(name, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *State) TypeName() string { return packageName }
