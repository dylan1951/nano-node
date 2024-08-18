package ledger

import (
	"github.com/accept-nano/ed25519-blake2b"
	"log"
	"node/blocks"
	"node/config"
	"node/store"
	"node/types"
)

var queue = make(chan blocks.Block, 1000)
var blocked = make(map[types.Hash]blocks.Block)

func init() {
	if store.GetBlock(config.Network.Genesis.Hash()) == nil {
		println("Initializing ledger with genesis block")
		store.PutBlock(&config.Network.Genesis)
		store.SetAccount(config.Network.Genesis.Account, store.Account{Frontier: config.Network.Genesis.Hash()})
	}
}

func ProcessBlocks() {
	for block := range queue {
		hash := block.Hash()
		switch block := block.(type) {
		case *blocks.OpenBlock:
			if !ed25519.Verify(block.Account[:], hash[:], block.Signature[:]) {
				log.Fatalf("block has invalid signature")
			}
			if store.GetAccount(block.Account) != nil {
				log.Fatalf("open block when account already exists")
			}
		case *blocks.SendBlock:

		}

		store.PutBlock(block)
		store.SetAccount(block.Account, store.Account{Frontier: block.Hash()})
	}
}
