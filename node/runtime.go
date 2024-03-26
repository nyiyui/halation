package node

import (
	"log"
	"sync"

	"nyiyui.ca/halation/aiz"
)

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

func (nr *NodeRunner) ActivateNodeUsingPromises(nn NodeName, doneCh chan<- struct{}) {
	var node Node
	func() {
		nr.NMLock.RLock()
		defer nr.NMLock.RUnlock()
		var ok bool
		node, ok = nr.NM.Nodes[nn]
		if !ok {
			panic("node not found by name")
		}
		node = node.Clone()
	}()
	go func() {
		result, err := node.Activate(nr.runner, nil)
		if err != nil {
			log.Printf("activating node: %s", err)
		}
		var wg sync.WaitGroup
		pm := nr.NM.GenPromiseMap()
		for _, user := range pm[nn] {
			func() {
				nr.NMLock.Lock()
				defer nr.NMLock.Unlock()
				for _, promise := range nr.NM.Nodes[user].BaseNodeRef().Promises {
					if promise.SupplyNodeName != nn {
						continue
					}
					if promise.FieldName != "dummy" {
						setValue(nr.NM.Nodes[user], promise.FieldName, result.String())
					}
				}
			}()
			innerDoneCh := make(chan struct{})
			nr.ActivateNodeUsingPromises(user, innerDoneCh)
			wg.Add(1)
			go func() {
				<-innerDoneCh
				wg.Done()
			}()
		}
		if doneCh != nil {
			wg.Wait()
			doneCh <- struct{}{}
		}
	}()
}
