package node

import (
	"fmt"
	"time"

	"nyiyui.ca/halation/aiz"
)

type Timer struct {
	*BaseNode
	Delay time.Duration `json:"delay"`
}

func newTimerBlank() *Timer {
	return &Timer{BaseNode: new(BaseNode)}
}

func NewTimer(delay time.Duration) *Timer {
	return &Timer{BaseNode: new(BaseNode), Delay: delay}
}

func (t *Timer) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	<-time.After(t.Delay)
	return nil, nil
}

func (t *Timer) TypeName() string { return "nyiyui.ca/halation/node.Timer" }
