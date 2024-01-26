package node

import (
	"log"
	"sync"
)

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
			return
		}
		var wg sync.WaitGroup
		pm := nr.NM.genPromiseMap()
		for _, user := range pm[nn] {
			func() {
				nr.NMLock.Lock()
				defer nr.NMLock.Unlock()
				for _, promise := range nr.NM.Nodes[user].BaseNodeRef().Promises {
					if promise.SupplyNodeName != nn {
						continue
					}
					setValue(nr.NM.Nodes[user], promise.FieldName, result.String())
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
