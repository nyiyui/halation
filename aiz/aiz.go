// Package aiz contains general cue structures.
package aiz

import "nyiyui.ca/halation/osc"

type Runner struct {
	OSCClient *osc.Client
}

type Show struct {
	Cues []Cue `json:"cues"`
}

type Cue struct {
	SGs             []SG     `json:"sgs"`
	DefaultGradient Gradient `json:"gradient"`
}

type SG struct {
	State    State    `json:"state"`
	Gradient Gradient `json:"gradient"`
}

type State interface {
	Reify(r *Runner, g Gradient)
}

var _ = []State{
	new(osc.State),
}

// Gradient provides a transition from different states.
// All integers values (unless specified otherwise) are in microseconds.
// Values are from 0 to 1 inclusive.
type Gradient interface {
	Duration() int
	PreferredResolution() int
	ValueAt(t int) float32
	Values(resolution int) []float32
}

var _ = []Gradient{
	new(LinearGradient),
}
