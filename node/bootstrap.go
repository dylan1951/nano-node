package node

import (
	"fmt"
)

func Bootstrap() {
	// scan frontiersChan
	start := [32]byte{}

	for {
		for _, peer := range peers {
			if peer != nil {
				//fmt.Printf("requesting 1000 frontiers from %s\n", peer.AddrPort().Addr().String())
				frontiers := <-peer.RequestFrontiers(start, 1000)
				fmt.Printf("recieved %d frontiers starting at %s, from %s\n", len(frontiers), frontiers[0].Account.GoString(), peer.AddrPort().Addr().String())
				start = frontiers[len(frontiers)-1].Account
			}
		}
	}
}
