package node

import (
	"encoding/json"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
	"nyiyui.ca/halation/aiz"
)

type EvalLua struct {
	*BaseNode
	Source string `halation:""`
}

func newEvalLuaBlank() *EvalLua {
	s := new(EvalLua)
	s.BaseNode = new(BaseNode)
	return s
}

func (s *EvalLua) Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	l := lua.NewState()
	defer l.Close()
	luajson.Preload(l)
	l.SetGlobal("params", lua.LString(data))
	if err := l.DoString(s.Source); err != nil {
		return nil, err
	}
	ret := l.Get(-1)
	return ret, nil
}

func (s *EvalLua) Clone() Node {
	return &EvalLua{
		BaseNode: s.CloneBaseNode(),
		Source:   s.Source,
	}
}

func (s *EvalLua) TypeName() string { return "nyiyui.ca/halation/node.EvalLua" }
