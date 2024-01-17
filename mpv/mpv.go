// Package mpv controls mpv to play media.
package mpv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/blang/mpv"
	"nyiyui.ca/halation/aiz"
)

const packageName = "nyiyui.ca/halation/mpv"

func init() {
	aiz.StateTypes[packageName] = func() aiz.State { return new(State) }
}

type Client struct {
	client *mpv.Client
}

func NewClient(client *mpv.Client) *Client {
	return &Client{client}
}

func NewClientUsingSubprocess() (c *Client, err error) {
	var f *os.File
	f, err = os.CreateTemp("", "halation-mpv-ipc-*")
	if err != nil {
		return
	}
	ipcPath := f.Name()
	err = f.Close()
	if err != nil {
		return
	}
	cmd := exec.Command("mpv", "--idle", fmt.Sprintf("--input-ipc-server=%s", ipcPath))
	err = cmd.Start()
	if err != nil {
		return
	}
	time.Sleep(1 * time.Second)
	c = &Client{mpv.NewClient(mpv.NewIPCClient(ipcPath))}
	return
}

func (c *Client) Register(r *aiz.Runner) {
	r.Specific[packageName] = c
}

type State struct {
	FilePath string
	Paused   *bool
	// TODO: Position doesn't work when the file is initially loaded.
	Position        *int
	Fullscreen      *bool
	ExtraProperties map[string]interface{} `json:"extraProperties"`
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
