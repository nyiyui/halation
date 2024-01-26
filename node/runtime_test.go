package node

import (
	"testing"

	"nyiyui.ca/halation/aiz"
)

func TestRuntime(t *testing.T) {
	nr := NewNodeRunner(new(aiz.Runner))
	a := NodeName{"", "a"}
	b := NodeName{"", "b"}
	a2 := newEvalLuaBlank()
	a2.Source = `print("a")
return "print(\"b new\")"`
	nr.NM.Nodes[a] = a2
	b2 := newEvalLuaBlank()
	b2.Source = `print("b original")`
	b2.BaseNodeRef().Promises = append(b2.BaseNodeRef().Promises, Promise{
		FieldName:      "Source",
		SupplyNodeName: a,
	})
	nr.NM.Nodes[b] = b2
	doneCh := make(chan struct{})
	nr.ActivateNodeUsingPromises(a, doneCh)
	<-doneCh
	// TODO: check that the code running was "b new"
}
