package aiz

import (
	"context"

	"nyiyui.ca/halation/notify"
)

type SeriesName = string

type Runner struct {
	Specific      map[string]interface{}
	CurrentStates map[SeriesName]map[string]State

	sgAppliedMuxS *notify.MultiplexerSender[*SG]
	SGAppliedMux  *notify.Multiplexer[*SG]
}

func NewRunner() *Runner {
	r := new(Runner)
	r.Specific = map[string]interface{}{}
	r.CurrentStates = map[SeriesName]map[string]State{}
	r.sgAppliedMuxS, r.SGAppliedMux = notify.NewMultiplexerSender[*SG]("Runner.sgAppliedMuxS")
	return r
}

// ApplySG applies the given sg.
// sg must not be mutated once this function is called.
func (r *Runner) ApplySG(sg *SG, ctx context.Context) error {
	r.sgAppliedMuxS.Send(sg)
	_, ok := r.CurrentStates[sg.Series]
	if !ok {
		r.CurrentStates[sg.Series] = map[string]State{}
	}
	prev := r.CurrentStates[sg.Series][sg.State.TypeName()]
	defer func() {
		r.CurrentStates[sg.Series][sg.State.TypeName()] = sg.State
	}()
	return sg.State.Reify(r, sg.Gradient, prev)
}
