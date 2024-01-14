// Package aiz contains general cue structures.
package aiz

import (
	"encoding/json"
	"fmt"
)

type Runner struct {
	Specific map[string]interface{}
}

type Show struct {
	Cues []Cue `json:"cues"`
}

type Cue struct {
	SGs []SG `json:"sgs"`
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
	stateJ, err := json.Marshal(sg.State)
	if err != nil {
		return nil, fmt.Errorf("marshal State: %w", err)
	}
	gradientJ, err := json.Marshal(sg.Gradient)
	if err != nil {
		return nil, fmt.Errorf("marshal Gradient: %w", err)
	}
	j := sgJSON{
		StateType:    sg.State.TypeName(),
		State:        stateJ,
		GradientType: sg.Gradient.TypeName(),
		Gradient:     gradientJ,
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
