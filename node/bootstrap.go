package node

import (
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	"node/blocks"
	"node/ledger"
	"node/messages"
	"node/store"
)

const pullSize = 128

func Bootstrap() {
	println("starting bootstrap")

	account := ledger.GetUnsyncedAccount()

	for {
		for _, peer := range peers {
			if peer == nil {
				continue
			}

			var pull []blocks.Block
			var ok bool

			if account.Height == 0 {
				fmt.Printf("asking for %d blocks starting from account %v from %s\n", pullSize, account.PublicKey.GoString(), peer.AddrPort().String())
				pull, ok = <-peer.RequestBlocks(account.PublicKey, pullSize, messages.Account)
			} else {
				fmt.Printf("asking for %d blocks starting from block %v from %s\n", pullSize, account.Frontier.GoString(), peer.AddrPort().String())
				pull, ok = <-peer.RequestBlocks(account.Frontier, pullSize, messages.Block)
			}

			if !ok {
				panic("peer died")
			}

			fmt.Printf("received %d blocks from %s\n", len(pull), peer.AddrPort().String())

			batch := store.NewBatch()

			missingDep := processBlocks(batch, account, pull)

			if err := batch.Commit(nil); err != nil {
				panic(err)
			}

			if missingDep != nil {
				account = ledger.GetUnsyncedAccount()
			} else if len(pull) < pullSize {
				fmt.Printf("account %v is fully synced at frontier: %v\n", account.PublicKey.GoString(), account.Frontier.GoString())
				store.MarkAccountSynced(account.PublicKey)
				account = ledger.GetUnsyncedAccount()
			}
		}
	}
}

func processBlocks(batch *pebble.Batch, account *ledger.Account, blocks []blocks.Block) error {
	for _, block := range blocks {
		if err := account.AddBlock(batch, block); err != nil {
			var missingDep *ledger.MissingDependency
			switch {
			case errors.Is(err, ledger.Invalid):
				panic("block is invalid")
			case errors.Is(err, ledger.Fork):
				panic("block is a fork")
			case errors.Is(err, ledger.Old):
				// ignore for now
			case errors.As(err, &missingDep):
				fmt.Printf("missing dependency %v\n", missingDep.Dependency.GoString())
				store.MarkAccountBlocked(batch, account.PublicKey, missingDep.Dependency)
				return err
			default:
				panic(err)
			}
		}
		store.MarkDependencyResolved(batch, block.Hash())
	}
	return nil
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
