package node

import (
	"fmt"
	"log"
	"sync"

	"nyiyui.ca/halation/aiz"
)

type NodeName struct {
	Package string
	Name    string
}

type NodeRequest struct {
	Params fmt.Stringer
}

type Node interface {
	GetListensTo() []NodeName
	// Activate can block.
	Activate(r *aiz.Runner, params fmt.Stringer) (result fmt.Stringer, err error)
	TypeName() string
}

type BaseNode struct {
	ListensTo []NodeName
}

func (b *BaseNode) GetListensTo() []NodeName { return b.ListensTo }

type NodeMap struct {
	Nodes map[NodeName]Node
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
	nm     *NodeMap
	nmLock sync.RWMutex
}

func (nr *NodeRunner) ActivateNode(nn NodeName, params fmt.Stringer) {
	var node Node
	func() {
		nr.nmLock.RLock()
		defer nr.nmLock.RUnlock()
		var ok bool
		node, ok = nr.nm.Nodes[nn]
		if !ok {
			panic("node not found by name")
		}
	}()
	go func() {
		result, err := node.Activate(nr.runner, params)
		if err != nil {
			log.Printf("activating node: %s", err)
			return
		}
		nr.nmLock.RLock()
		defer nr.nmLock.RUnlock()
		listeners := nr.nm.genListeners()
		for _, listener := range listeners[nn] {
			nr.ActivateNode(listener, result)
		}
	}()
}

var _ = []Node{
	new(EvalLua),
	new(Manual),
	new(SetState),
	new(Timer),
}
