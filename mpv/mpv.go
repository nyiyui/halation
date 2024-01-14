// Package mpv controls mpv to play media.
package mpv

import (
	"errors"
	"fmt"
	"github.com/blang/mpv"
	"log"
	"nyiyui.ca/halation/aiz"
	"os"
	"os/exec"
	"time"
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
	Paused   bool
	// TODO: Position doesn't work when the file is initially loaded.
	Position        int
	Fullscreen      bool
	ExtraProperties map[string]interface{}
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
	if currentPath != s.FilePath {
		err = c.client.Loadfile(s.FilePath, mpv.LoadFileModeReplace)
		if err != nil {
			return err
		}
	}

	err = c.client.SetPause(s.Paused)
	if err != nil {
		return err
	}

	log.Printf("seeking to %d", s.Position)
	err = c.client.Seek(s.Position, mpv.SeekModeAbsolute)
	if err != nil {
		return err
	}

	err = c.client.SetFullscreen(s.Fullscreen)
	if err != nil {
		return err
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
