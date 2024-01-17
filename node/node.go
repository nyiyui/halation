package node

import (
	"fmt"
	"log"
	"strings"
	"sync"

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

type NodeRequest struct {
	Params fmt.Stringer
}

// Node is not goroutine-safe.
type Node interface {
	GetDescription() string
	SetDescription(string)
	GetListensTo() []NodeName
	SetListensTo([]NodeName)
	// Activate can block.
	Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error)
	TypeName() string
}

type BaseNode struct {
	Description string
	ListensTo   []NodeName
}

func (b *BaseNode) GetDescription() string { return b.Description }

func (b *BaseNode) SetDescription(d string) { b.Description = d }

func (b *BaseNode) GetListensTo() []NodeName { return b.ListensTo }

func (b *BaseNode) SetListensTo(listensTo []NodeName) { b.ListensTo = listensTo }

type NodeMap struct {
	Nodes map[NodeName]Node
}

func NewNodeMap() *NodeMap {
	return &NodeMap{
		Nodes: map[NodeName]Node{},
	}
}

func (nm *NodeMap) GenListeners() map[NodeName][]NodeName {
	return nm.genListeners()
}

func (nm *NodeMap) genListeners() map[NodeName][]NodeName {
	listeners := map[NodeName][]NodeName{}
	for listener, node := range nm.Nodes {
		for _, listenee := range node.GetListensTo() {
			listeners[listenee] = append(listeners[listenee], listener)
		}
	}
	return listeners
}

type NodeRunner struct {
	runner *aiz.Runner
	NM     *NodeMap
	NMLock sync.RWMutex
}

func NewNodeRunner(runner *aiz.Runner) *NodeRunner {
	return &NodeRunner{
		runner: runner,
		NM:     NewNodeMap(),
	}
}

func (nr *NodeRunner) ActivateNode(nn NodeName, params fmt.Stringer) {
	var node Node
	func() {
		nr.NMLock.RLock()
		defer nr.NMLock.RUnlock()
		var ok bool
		node, ok = nr.NM.Nodes[nn]
		if !ok {
			panic("node not found by name")
		}
	}()
	go func() {
		nr.NMLock.RLock()
		defer nr.NMLock.RUnlock()
		result, err := node.Activate(nr.runner, params)
		if err != nil {
			log.Printf("activating node: %s", err)
			return
		}
		nr.NMLock.RLock()
		defer nr.NMLock.RUnlock()
		listeners := nr.NM.genListeners()
		for _, listener := range listeners[nn] {
			nr.ActivateNode(listener, result)
		}
	}()
}
