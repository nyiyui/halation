package node

import (
	"fmt"
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/timeutil"
)

type Timer struct {
	*BaseNode
	Delay timeutil.Duration `json:"delay"`
}

func newTimerBlank() *Timer {
	return &Timer{BaseNode: new(BaseNode)}
}

func NewTimer(delay time.Duration) *Timer {
	return &Timer{BaseNode: new(BaseNode), Delay: timeutil.Duration(delay)}
}

func (t *Timer) Clone() Node {
	return &Timer{
		BaseNode: t.CloneBaseNode(),
		Delay:    t.Delay,
	}
}

func (t *Timer) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	<-time.After(time.Duration(t.Delay))
	return nil, nil
}

func (t *Timer) TypeName() string { return "nyiyui.ca/halation/node.Timer" }
