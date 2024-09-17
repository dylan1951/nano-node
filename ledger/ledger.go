package ledger

import (
	"node/config"
	"node/store"
	"node/types/uint128"
)

func Init() {
	if store.GetBlock(config.Network.Genesis.Hash()) == nil {
		println("Initializing ledger with genesis block")
		store.PutBlock(config.Network.Genesis.Hash(), store.BlockRecord{
			Block:   &config.Network.Genesis,
			Account: config.Network.Genesis.Account,
		})
		store.SetAccount(config.Network.Genesis.Account, store.AccountRecord{
			Frontier: config.Network.Genesis.Hash(),
			Height:   1,
			Balance:  uint128.Max,
		})
		store.MarkAccountUnsynced(config.Network.Genesis.Account)
	}
}
