package ledger

import (
	"errors"
	"fmt"
	"node/blocks"
	"node/config"
	"node/store"
	"node/types"
	"node/types/uint128"
)

var Queue = make(chan blocks.Block, 1000)
var blocked = make(map[types.Hash]blocks.Block)

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
	}
}

func ProcessBlocks() {
	for block := range Queue {
		publicKey := PubKeyFromBlock(block)

		if publicKey == nil {
			fmt.Printf("missing dependency %#v\n", block.GetPrevious())
			continue
		}

		fmt.Printf("ledger processing block: %s\n", block.Hash().GoString())

		account := AccountFromPublicKey(*publicKey)
		err := account.AddBlock(block)

		var missingDep MissingDependency
		switch {
		case errors.Is(err, Invalid):
			println("block is invalid")
		case errors.Is(err, Fork):
			println("block is a fork")
		case errors.As(err, &missingDep):
			fmt.Printf("missing dependency %#v\n", missingDep.Dependency)
			blocked[missingDep.Dependency] = block
		}
	}
}
