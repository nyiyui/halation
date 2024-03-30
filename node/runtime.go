package node

import (
	"fmt"
	"log"
	"sync"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/notify"
)

type Change struct {
	NodeName  NodeName
	Activated bool
	Error     string
}

type NodeRunner struct {
	runner     *aiz.Runner
	NM         *NodeMap
	NMLock     sync.RWMutex
	changeMuxS *notify.MultiplexerSender[Change]
}

func NewNodeRunner(runner *aiz.Runner) *NodeRunner {
	return &NodeRunner{
		runner: runner,
		NM:     NewNodeMap(),
	}
}

func (nr *NodeRunner) SetChangeMuxS(changeMuxS *notify.MultiplexerSender[Change]) {
	nr.changeMuxS = changeMuxS
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
		result, activateErr := node.Activate(nr.runner, nil)
		if activateErr != nil {
			log.Printf("activating node: %s", activateErr)
		}
		if nr.changeMuxS != nil {
			msg := fmt.Sprint(activateErr)
			if activateErr == nil {
				msg = ""
			}
			nr.changeMuxS.Send(Change{NodeName: nn, Activated: true, Error: msg})
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
