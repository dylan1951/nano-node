package node

import (
	"fmt"
	"node/config"
	"node/ledger"
	"node/types"
)

var pool = make(chan types.Hash, 1000)

func Bootstrap() {
	println("starting bootstrap")

	pool <- config.Network.Genesis.Hash()

	for {
		for _, peer := range peers {
			if peer != nil {
				hash := <-pool
				fmt.Printf("requesting 10 blocks starting from %s from %s\n", hash.GoString(), peer.AddrPort().String())
				blocks := <-peer.RequestBlocks(hash, 10)
				if len(blocks) != 10 {
					fmt.Printf("received %d blocks instead of 10\n", len(blocks))
				} else {
					fmt.Printf("received %d blocks, adding the latest hash back to the pool\n", len(blocks))
					pool <- blocks[9].Hash()
				}
				for _, block := range blocks {
					ledger.Queue <- block
				}
			}
		}
	}
}

func ScanFrontiers() {
	start := [32]byte{}

	for {
		for _, peer := range peers {
			if peer != nil {
				frontiers := <-peer.RequestFrontiers(start, 1000)
				fmt.Printf("recieved %d frontiers starting at %s, from %s\n", len(frontiers), frontiers[0].Account.GoString(), peer.AddrPort().Addr().String())
				start = frontiers[len(frontiers)-1].Account
			}
		}
	}
}
