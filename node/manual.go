package node

import (
	"fmt"

	"nyiyui.ca/halation/aiz"
)

type Manual struct {
	*BaseNode
}

func NewManual() *Manual {
	m := new(Manual)
	m.BaseNode = new(BaseNode)
	return m
}

func (m *Manual) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	return nil, nil
}

func (m *Manual) Clone() Node {
	return &Manual{
		BaseNode: m.CloneBaseNode(),
	}
}

func (m *Manual) TypeName() string { return "nyiyui.ca/halation/node.Manual" }
