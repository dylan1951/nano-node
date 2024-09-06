package store

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
	"node/types"
	"node/types/uint128"
	"node/utils"
)

var db, _ = pebble.Open("data", &pebble.Options{})

const (
	PrefixBlock byte = iota
	PrefixAccount
)

type BlockRecord struct {
	Block      blocks.Block
	Account    types.PublicKey
	Receivable uint128.Uint128
}

type AccountRecord struct {
	Frontier types.Hash
	Balance  uint128.Uint128
	Height   uint64
}

func (b *BlockRecord) Serialize() []byte {
	return []byte{}
}

func (b *BlockRecord) Deserialize([]byte) {

}

func (a *AccountRecord) Serialize() []byte {
	return []byte{}
}

func (a *AccountRecord) Deserialize([]byte) {

}

func SetAccount(publicKey [32]byte, account AccountRecord) {
	if err := db.Set(publicKey[:], utils.Serialize(account, binary.LittleEndian), pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func GetAccount(publicKey types.PublicKey) *AccountRecord {
	serialized, closer, err := db.Get(publicKey[:])

	if err != nil {
		return nil
	}

	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}

	return utils.Read[AccountRecord](bytes.NewReader(serialized), binary.LittleEndian)
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

func PutBlock(blockHash types.Hash, record BlockRecord) {
	if err := db.Set(blockHash[:], record.Serialize(), pebble.Sync); err != nil {
		log.Fatal(err)
	}

	fmt.Println("saved block:", hex.EncodeToString(blockHash[:]))
}

func GetBlock(blockHash [32]byte) *BlockRecord {
	serialized, closer, err := db.Get(blockHash[:])

	if err != nil {
		return nil
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
		return [32]byte(iter.Key())
	} else {
		return [32]byte{}
	}
}
