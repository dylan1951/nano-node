package store

import (
	"encoding/hex"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
)

var db, _ = pebble.Open("data", &pebble.Options{})

const (
	PrefixBlock byte = iota
	PrefixAccount
)

type Account struct {
	Frontier [32]byte
}

func SetPull(hash [32]byte, val byte) {
	if err := db.Set(hash[:], []byte{val}, pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func DeletePull(hash [32]byte) {
	if err := db.Delete(hash[:], pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func PutBlock(block blocks.Block) {
	blockHash := block.Hash()

	if err := db.Set(blockHash[:], block.Serialize(), pebble.Sync); err != nil {
		log.Fatal(err)
	}

	fmt.Println("saved block:", hex.EncodeToString(blockHash[:]))
}

func GetBlock(blockHash [32]byte) blocks.Block {
	serialized, closer, err := db.Get(blockHash[:])

	if err != nil {
		log.Fatal(err)
	}

	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}

	return blocks.Deserialize(serialized)
}

func GetLastBlockHash() [32]byte {
	iter, err := db.NewIter(&pebble.IterOptions{})
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
