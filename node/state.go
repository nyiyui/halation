package node

import (
	"context"
	"fmt"

	"nyiyui.ca/halation/aiz"
)

type SetState struct {
	*BaseNode
	Runner *aiz.Runner
	SG     *aiz.SG
}

func NewSetState() *SetState {
	s := new(SetState)
	return s
}

func (s *SetState) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	return nil, s.Runner.ApplySG(s.SG, context.Background())
}

func (s *SetState) TypeName() string { return "nyiyui.ca/halation/node.SetState" }
