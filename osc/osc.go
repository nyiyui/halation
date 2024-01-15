// Package osc provides an interface to the ColorSource AV lightboard using OSC commands over HTTP.
package osc

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"nyiyui.ca/halation/aiz"
)

const packageName = "nyiyui.ca/halation/osc"

func init() {
	aiz.StateTypes["nyiyui.ca/halation/osc"] = func() aiz.State { return new(State) }
}

const NumberOfChannels = 40

type Client struct {
	HTTP    *http.Client
	BaseURL *url.URL
}

func NewDefaultClient() *Client {
	u, err := url.Parse("http://10.10.0.2:8080")
	if err != nil {
		panic(err)
	}
	return &Client{
		HTTP: &http.Client{
			Timeout: 1 * time.Second,
		},
		BaseURL: u,
	}
}

func (c *Client) Register(r *aiz.Runner) {
	r.Specific[packageName] = c
}

func (c *Client) Blackout() (err error) {
	err = c.ChanSelect(1)
	if err != nil {
		return
	}
	err = c.ChanThru(40)
	if err != nil {
		return
	}
	err = c.ChanAt(0)
	if err != nil {
		return
	}
	return
}

func (c *Client) sendGet(command string) error {
	r, err := c.HTTP.Get(c.BaseURL.JoinPath(command).String())
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("returned non-200 code: %d", r.StatusCode)
	}
	return nil
}

// ChanSelect selects the given channel.
func (c *Client) ChanSelect(channel int) error {
	return c.sendGet(fmt.Sprintf("cs/chan/select/%d", channel))
}

// ChanAdd adds the channel to the selection.
func (c *Client) ChanAdd(channel int) error {
	return c.sendGet(fmt.Sprintf("cs/chan/add/%d", channel))
}

// ChanSubtract adds the channel to the selection.
func (c *Client) ChanSubtract(channel int) error {
	return c.sendGet(fmt.Sprintf("cs/chan/subtract/%d", channel))
}

// ChanThru selects every channel form the last channel selected thru to the given channel.
func (c *Client) ChanThru(channel int) error {
	return c.sendGet(fmt.Sprintf("cs/chan/thru/%d", channel))
}

// ChanSet selects listed channels and them only.
func (c *Client) ChanSet(channels []int) error {
	if len(channels) == 0 {
		err := c.ChanSelect(1)
		if err != nil {
			return err
		}
		return c.ChanSubtract(1)
	}
	err := c.ChanSelect(channels[0])
	if err != nil {
		return err
	}
	for i := range channels {
		if i == 0 {
			continue
		}
		err = c.ChanAdd(channels[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// ChanAt sets the selected channel's level.
func (c *Client) ChanAt(level int) error {
	return c.sendGet(fmt.Sprintf("cs/chan/at/%d", level))
}

// ColorHS sets the selected channel's color.
// Hue is in the inclusive range of [0, 360].
// Saturation is in the inclusive range of [0, 100].
func (c *Client) ColorHS(hue, saturation int) error {
	return c.sendGet(fmt.Sprintf("cs/color/hs/%d/%d", hue, saturation))
}
