package node

import (
	"fmt"
	"strings"

	"nyiyui.ca/halation/aiz"
)

var NodeTypes = map[string]func() Node{}

var newNodeFns = []func() Node{
	func() Node { return newEvalLuaBlank() },
	func() Node { return NewManual() },
	func() Node { return newSetStateBlank() },
	func() Node { return newTimerBlank() },
}

func init() {
	for _, newNode := range newNodeFns {
		n := newNode()
		NodeTypes[n.TypeName()] = newNode
	}
}

type NodeName struct {
	Package string
	Name    string
}

func ParseNodeName(s string) NodeName {
	i := strings.LastIndex(s, ".")
	if i == -1 {
		return NodeName{"", s}
	}
	return NodeName{s[:i], s[i+1:]}
}

func (nn NodeName) String() string {
	return nn.Package + "." + nn.Name
}

func (nn NodeName) IsZero() bool {
	return nn.Package == "" && nn.Name == ""
}

type NodeRequest struct {
	Params fmt.Stringer
}

// Node is not goroutine-safe.
// Nodes are the basic building blocks for the Halation runtime. See node/README.md for details.
type Node interface {
	BaseNodeRef() *BaseNode
	Clone() Node
	GetDescription() string
	SetDescription(string)
	// Activate can block.
	Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error)
	TypeName() string
}

type BaseNode struct {
	Description string
	Promises    []Promise
}

func (b *BaseNode) CloneBaseNode() *BaseNode {
	b2 := &BaseNode{
		Description: b.Description,
		Promises:    make([]Promise, len(b.Promises)),
	}
	copy(b2.Promises, b.Promises)
	return b2
}

func (b *BaseNode) BaseNodeRef() *BaseNode { return b }

func (b *BaseNode) GetDescription() string { return b.Description }

func (b *BaseNode) SetDescription(d string) { b.Description = d }

type NodeMap struct {
	Nodes map[NodeName]Node
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		Nodes: map[NodeName]Node{},
	}
}

func (nm *NodeMap) GenPromiseMap() map[NodeName][]NodeName {
	pm := map[NodeName][]NodeName{}
	for user, node := range nm.Nodes {
		for _, promise := range node.BaseNodeRef().Promises {
			pm[promise.SupplyNodeName] = append(pm[promise.SupplyNodeName], user)
		}
	}
	return pm
}
