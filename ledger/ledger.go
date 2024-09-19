package ledger

import (
	"node/config"
	"node/store"
	"node/types/uint128"
)

func Init() {
	if store.GetBlock(config.Network.Genesis.Hash()) == nil {
		println("Initializing ledger with genesis block")
		batch := store.NewBatch()

		store.PutBlock(batch, config.Network.Genesis.Hash(), store.BlockRecord{
			Block:   &config.Network.Genesis,
			Account: config.Network.Genesis.Account,
		})
		store.SetAccount(batch, config.Network.Genesis.Account, store.AccountRecord{
			Frontier: config.Network.Genesis.Hash(),
			Height:   1,
			Balance:  uint128.Max,
		})

		store.MarkAccountUnsynced(batch, config.Network.Genesis.Account)

		if err := batch.Commit(nil); err != nil {
			panic(err)
		}
	}
}
