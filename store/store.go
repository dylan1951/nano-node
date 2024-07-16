package store

import (
	"encoding/hex"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
)

var blockDB, _ = pebble.Open("data/blocks", &pebble.Options{})
var accountDB, _ = pebble.Open("data/accounts", &pebble.Options{})
var pullsDB, _ = pebble.Open("data/pulls", &pebble.Options{})

func SetPull(hash [32]byte, val byte) {
	if err := pullsDB.Set(hash[:], []byte{val}, pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func DeletePull(hash [32]byte) {
	if err := pullsDB.Delete(hash[:], pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func PutBlock(block blocks.Block) {
	blockHash := block.Hash()

	if err := blockDB.Set(blockHash[:], block.Serialize(), pebble.Sync); err != nil {
		log.Fatal(err)
	}

	fmt.Println("saved block:", hex.EncodeToString(blockHash[:]))
}

func GetBlock(blockHash [32]byte) blocks.Block {
	serialized, closer, err := blockDB.Get(blockHash[:])

	if err != nil {
		log.Fatal(err)
	}

	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}

	return blocks.Deserialize(serialized)
}

func GetLastBlockHash() [32]byte {
	iter, err := blockDB.NewIter(nil)
	defer iter.Close()

	if err != nil {
		log.Fatal(err)
	}

	if iter.Last() {
		return *(*[32]byte)(iter.Key())
	} else {
		return [32]byte{}
	}
}
