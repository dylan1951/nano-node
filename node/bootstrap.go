package node

import (
	"node/types"
)

var pool []types.Hash
var accountRangesScanned []types.Hash

func (n *Node) Bootstrap() {
	for {
		for _, p := range n.peers {
			//telemetry, ok := <-p.RequestTelemetry()
			//if !ok {
			//	log.Fatal("Failed to get telemetry from peer")
			//}

			return
		}
	}
}

func (n *Node) scanAccounts() {

}
