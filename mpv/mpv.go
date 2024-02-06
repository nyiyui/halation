// Package mpv controls mpv to play media.
package mpv

import (
	"fmt"
	"os"
	"os/exec"
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
