package node

import (
	"fmt"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/notify"
)

type Manual struct {
	*BaseNode
	mux  *notify.Multiplexer[fmt.Stringer]
	MuxS *notify.MultiplexerSender[fmt.Stringer]
}

func NewManual() *Manual {
	m := new(Manual)
	m.MuxS, m.mux = notify.NewMultiplexerSender[fmt.Stringer]("")
	return m
}

func (m *Manual) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	return nil, nil
}

func (m *Manual) TypeName() string { return "nyiyui.ca/halation/node.Manual" }
