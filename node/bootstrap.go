package node

import (
	"encoding/hex"
	"fmt"
	"node/store"
	"node/types"
	"time"
)

var outdatedAccounts []types.PublicKey

func (n *Node) Bootstrap() {
	for _, p := range n.peers {
		start := [32]byte{}

		for {
			fmt.Printf("requesting frontiers from: %+v\n", hex.EncodeToString(start[:]))
			frontiers := <-p.RequestFrontiers(start)

			fmt.Printf("got %v frontiers\n", len(frontiers))

			if len(frontiers) == 0 {
				println("finished!")
				break
			}

			for _, frontier := range frontiers {
				//fmt.Printf("%s:%s\n", frontier.Account.GoString(), frontier.Hash.GoString())
				account := store.GetAccount(frontier.Account)
				if account == nil || account.Frontier != frontier.Account {
					store.SetAccount(frontier.Account, store.Account{Frontier: frontier.Hash})
					outdatedAccounts = append(outdatedAccounts, frontier.Account)
				}

				start = frontier.Account
			}

			time.Sleep(500 * time.Millisecond)
		}

		fmt.Printf("there are %v outdated accounts\n", len(outdatedAccounts))

		return
	}
}

func (n *Node) scanAccounts() {

}
