package node

import "sort"

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
