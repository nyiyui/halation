package node

import (
	"context"
	"fmt"

	"nyiyui.ca/halation/aiz"
)

type SetState struct {
	*BaseNode
	SG *aiz.SG
}

func newSetStateBlank() *SetState {
	return &SetState{BaseNode: new(BaseNode)}
}

func NewSetState(sg *aiz.SG) *SetState {
	return &SetState{BaseNode: new(BaseNode), SG: sg}
}

func (s *SetState) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	return nil, r.ApplySG(s.SG, context.Background())
}

func (s *SetState) Clone() Node {
	return &SetState{
		BaseNode: s.CloneBaseNode(),
		SG:       s.SG.Clone(),
	}
}

func (s *SetState) TypeName() string { return "nyiyui.ca/halation/node.SetState" }

func (s *SetState) String() string {
	return fmt.Sprintf("%s with gradient %v", s.SG.State, s.SG.Gradient)
}
