package node

import (
	"encoding/json"
	"sort"
)

type Cuelist struct {
	Nodes map[float64]NodeName

	CurrentCue float64
}

func NewCuelist() *Cuelist {
	return &Cuelist{
		Nodes:      map[float64]NodeName{},
		CurrentCue: 0,
	}
}

type cuelistJSON struct {
	Cues       []cueJSON
	CurrentCue float64
}

type cueJSON struct {
	Number float64
	Name   NodeName
}

func (c *Cuelist) MarshalJSON() ([]byte, error) {
	clj := cuelistJSON{
		Cues:       make([]cueJSON, 0, len(c.Nodes)),
		CurrentCue: c.CurrentCue,
	}
	for number, nn := range c.Nodes {
		clj.Cues = append(clj.Cues, cueJSON{number, nn})
	}
	return json.Marshal(clj)
}

func (c *Cuelist) UnmarshalJSON(data []byte) error {
	var clj cuelistJSON
	err := json.Unmarshal(data, &clj)
	if err != nil {
		return err
	}
	c.CurrentCue = clj.CurrentCue
	for _, cj := range clj.Cues {
		c.Nodes[cj.Number] = cj.Name
	}
	return nil
}

func (c *Cuelist) genIndices() []float64 {
	indices := make([]float64, 0, len(c.Nodes))
	for key := range c.Nodes {
		indices = append(indices, key)
	}
	sort.Float64s(indices)
	return indices
}

func (c *Cuelist) GenOpposite() map[NodeName]float64 {
	m := map[NodeName]float64{}
	for key, val := range c.Nodes {
		m[val] = key
	}
	return m
}
