package store

import (
	"encoding/hex"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
	"node/utils"
)

var db, _ = pebble.Open("data", &pebble.Options{})

func SaveBlock(block blocks.Block) {
	serialized := utils.Serialize(block)
	blockHash := block.Hash()

	if err := db.Set(blockHash[:], serialized, pebble.Sync); err != nil {
		log.Fatal(err)
	}

	db.
		fmt.Println("saved block:", hex.EncodeToString(blockHash[:]))
}
