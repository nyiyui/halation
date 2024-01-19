package aiz

import (
	"context"

	"nyiyui.ca/halation/notify"
)

type Runner struct {
	Specific      map[string]interface{}
	CurrentStates map[string]State

	sgAppliedMuxS *notify.MultiplexerSender[*SG]
	SGAppliedMux  *notify.Multiplexer[*SG]
}

func NewRunner() *Runner {
	r := new(Runner)
	r.Specific = map[string]interface{}{}
	r.CurrentStates = map[string]State{}
	r.sgAppliedMuxS, r.SGAppliedMux = notify.NewMultiplexerSender[*SG]("Runner.sgAppliedMuxS")
	return r
}

// ApplySG applies the given sg.
// sg must not be mutated once this function is called.
func (r *Runner) ApplySG(sg *SG, ctx context.Context) error {
	r.sgAppliedMuxS.Send(sg)
	prev := r.CurrentStates[sg.State.TypeName()]
	defer func() {
		r.CurrentStates[sg.State.TypeName()] = sg.State
	}()
	return sg.State.Reify(r, sg.Gradient, prev)
}
