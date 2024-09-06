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
	Version  uint8
}

func (b *BlockRecord) Serialize() []byte {
	data := b.Block.Serialize()
	data = append(data, b.Account[:]...)
	data = append(data, b.Receivable.Bytes()...)
	return data
}

func (b *BlockRecord) Deserialize(data []byte) *BlockRecord {
	reader := bytes.NewReader(data)
	b.Block = blocks.Read(reader)
	_, _ = reader.Read(b.Account[:])
	b.Receivable = uint128.Read(reader)
	return b
}

func (a *AccountRecord) Serialize() []byte {
	data := a.Frontier[:]
	data = append(data, a.Balance.Bytes()...)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, a.Height)
	data = append(data, heightBytes...)
	return append(data, a.Version)
}

func (a *AccountRecord) Deserialize(data []byte) *AccountRecord {
	reader := bytes.NewReader(data)
	_, _ = reader.Read(a.Frontier[:])
	a.Balance = uint128.Read(reader)
	heightBytes := make([]byte, 8)
	_, _ = reader.Read(heightBytes)
	a.Height = binary.BigEndian.Uint64(heightBytes)
	a.Version, _ = reader.ReadByte()
	return a
}

func SetAccount(publicKey [32]byte, account AccountRecord) {
	if err := db.Set(publicKey[:], account.Serialize(), pebble.Sync); err != nil {
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

	return (&AccountRecord{}).Deserialize(serialized)
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

	return (&BlockRecord{}).Deserialize(serialized)
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
