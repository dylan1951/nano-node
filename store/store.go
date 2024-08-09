package store

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
	"node/utils"
)

var db, _ = pebble.Open("data", &pebble.Options{})

const (
	PrefixBlock byte = iota
	PrefixAccount
)

type Account struct {
	Frontier [32]byte
}

func SetAccount(publicKey [32]byte, account Account) {
	if err := db.Set(publicKey[:], utils.Serialize(account, binary.LittleEndian), pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func GetAccount(publicKey [32]byte) *Account {
	serialized, closer, err := db.Get(publicKey[:])

	if err != nil {
		return nil
	}

	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}

	return utils.Read[Account](bytes.NewReader(serialized), binary.LittleEndian)
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
