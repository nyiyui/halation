package node

import (
	"encoding/json"
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

func (nn NodeName) MarshalJSON() ([]byte, error) {
	return json.Marshal(nn.String())
}

func (nn *NodeName) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*nn = ParseNodeName(s)
	return nil
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

type nodeMapJSON struct {
	Nodes map[string]NodeJSON
}

type NodeJSON struct {
	NodeType string
	Node     Node
}

type nodeJSON2 struct {
	NodeType string
	Node     json.RawMessage
}

func (nj *NodeJSON) UnmarshalJSON(data []byte) error {
	var nj2 nodeJSON2
	err := json.Unmarshal(data, &nj2)
	if err != nil {
		return err
	}
	newNode, ok := NodeTypes[nj2.NodeType]
	if !ok {
		return fmt.Errorf("unsupported node type %s", nj2.NodeType)
	}
	n := newNode()
	err = json.Unmarshal(nj2.Node, n)
	if err != nil {
		return err
	}
	nj.NodeType = nj2.NodeType
	nj.Node = n
	return nil
}

func (nm *NodeMap) MarshalJSON() ([]byte, error) {
	nmj := nodeMapJSON{make(map[string]NodeJSON, len(nm.Nodes))}
	for nn, n := range nm.Nodes {
		nmj.Nodes[nn.String()] = NodeJSON{n.TypeName(), n}
	}
	return json.Marshal(nmj)
}

func (nm *NodeMap) UnmarshalJSON(data []byte) error {
	var nmj nodeMapJSON
	err := json.Unmarshal(data, &nmj)
	if err != nil {
		return err
	}
	for nns, n := range nmj.Nodes {
		nm.Nodes[ParseNodeName(nns)] = n.Node
	}
	return nil
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
