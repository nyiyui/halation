package trigger

import (
	"time"

	"nyiyui.ca/halation/aiz"
)

type Timed struct {
	Delay      time.Duration  `json:"delay"`
	CueRequest aiz.CueRequest `json:"cueRequest"`
}

func (t *Timed) Sendback(r *aiz.Runner, ch chan<- aiz.CueRequest) error {
	go func() {
		<-time.After(t.Delay)
		ch <- t.CueRequest
	}()
	return nil
}
