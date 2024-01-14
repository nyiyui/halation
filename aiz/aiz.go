// Package aiz contains general cue structures.
package aiz

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	Specific   map[string]interface{}
	CurrentCue int
}

type Show struct {
	Cues []Cue `json:"cues"`
}

func (s *Show) prevCueState(from int, stateType string) (i, j int, state State) {
	for i := from; i >= 0; i-- {
		cue := s.Cues[i]
		for j, sg := range cue.SGs {
			if sg.State.TypeName() == stateType {
				return i, j, sg.State
			}
		}
	}
	return -1, -1, nil
}

func (s *Show) ApplyCue(r *Runner, i int, ctx context.Context) error {
	cue := s.Cues[i]
	errs, ctx := errgroup.WithContext(ctx)
	for _, sg := range cue.SGs {
		errs.Go(func() error {
			var prev State
			if i != 0 {
				_, _, prev = s.prevCueState(i-1, sg.State.TypeName())
			}
			return sg.State.Reify(r, sg.Gradient, prev)
		})
	}
	r.CurrentCue = i
	return errs.Wait()
}

type Cue struct {
	Name string `json:"name"`
	SGs  []SG   `json:"sgs"`
}

type SG struct {
	State    State    `json:"state"`
	Gradient Gradient `json:"gradient"`
}

var StateTypes = map[string]func() State{}
var GradientTypes = map[string]func() Gradient{}

type sgJSON struct {
	StateType    string
	State        json.RawMessage
	GradientType string
	Gradient     json.RawMessage
}

func (sg *SG) MarshalJSON() ([]byte, error) {
	var j sgJSON
	var err error
	if sg.State != nil {
		j.State, err = json.Marshal(sg.State)
		if err != nil {
			return nil, fmt.Errorf("marshal State: %w", err)
		}
		j.StateType = sg.State.TypeName()
	}
	if sg.Gradient != nil {
		j.Gradient, err = json.Marshal(sg.Gradient)
		if err != nil {
			return nil, fmt.Errorf("marshal Gradient: %w", err)
		}
		j.GradientType = sg.Gradient.TypeName()
	}
	return json.Marshal(j)
}

func (sg *SG) UnmarshalJSON(data []byte) error {
	var j sgJSON
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	stateTypeNew, ok := StateTypes[j.StateType]
	if !ok {
		return fmt.Errorf("state: no type named %s", j.StateType)
	}
	s := stateTypeNew()
	err = json.Unmarshal(j.State, s)
	if !ok {
		return fmt.Errorf("unmarshal state: %s", err)
	}
	gradientTypeNew, ok := GradientTypes[j.GradientType]
	if !ok {
		return fmt.Errorf("gradient: no type named %s", j.GradientType)
	}
	g := gradientTypeNew()
	err = json.Unmarshal(j.Gradient, g)
	if !ok {
		return fmt.Errorf("unmarshal gradient: %s", err)
	}
	sg.State = s
	sg.Gradient = g
	return nil
}

type State interface {
	Reify(r *Runner, g Gradient, prev State) error
	TypeName() string
}

// Gradient provides a transition from different states.
// All integers values (unless specified otherwise) are in microseconds.
// Values are from 0 to 1 inclusive.
type Gradient interface {
	Duration() int
	PreferredResolution() int
	ValueAt(t int) float32
	Values(resolution int) []float32
	TypeName() string
}
