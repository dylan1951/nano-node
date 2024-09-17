package store

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	_ "github.com/kelindar/binary"
	"node/blocks"
	"node/types"
	"node/types/uint128"
	"node/utils"
)

var db, _ = pebble.Open("data", &pebble.Options{})

var (
	PrefixBlock           = []byte{0}
	PrefixAccount         = []byte{1}
	PrefixBlockedAccount  = []byte{2}
	PrefixUnsyncedAccount = []byte{3}
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

// MarkAccountUnsynced marks a new account as unsynced
func MarkAccountUnsynced(publicKey types.PublicKey) {
	fmt.Printf("Marking unsynced account %s\n", publicKey.GoString())
	key := append(PrefixUnsyncedAccount, publicKey[:]...)
	if err := db.Set(key, []byte{}, pebble.Sync); err != nil {
		panic(err)
	}
}

// MarkAccountSynced marks an account as synced
func MarkAccountSynced(publicKey types.PublicKey) {
	key := append(PrefixUnsyncedAccount, publicKey[:]...)
	if err := db.Delete(key, pebble.Sync); err != nil {
		panic(err)
	}
}

// MarkAccountBlocked marks an unsynced account as blocked by a given dependency.
func MarkAccountBlocked(publicKey types.PublicKey, dependency types.Hash) {
	key := append(PrefixUnsyncedAccount, publicKey[:]...)
	if err := db.Delete(key, pebble.Sync); err != nil {
		panic(err)
	}

	key = append(PrefixBlockedAccount, dependency[:]...)
	if err := db.Set(key, publicKey[:], pebble.Sync); err != nil {
		panic(err)
	}
}

// MarkAccountUnblocked marks an account as unblocked and moves it back to the unsynced set.
func MarkAccountUnblocked(publicKey types.PublicKey, dependency types.Hash) {
	key := append(PrefixBlockedAccount, dependency[:]...)
	if err := db.Delete(key, pebble.Sync); err != nil {
		panic(err)
	}

	key = append(PrefixUnsyncedAccount, dependency[:]...)
	if err := db.Set(key, publicKey[:], pebble.Sync); err != nil {
		panic(err)
	}
}

func GetUnsyncedAccount() (*types.PublicKey, error) {
	iter, _ := db.NewIter(&pebble.IterOptions{LowerBound: PrefixUnsyncedAccount})
	defer iter.Close()

	// Seek to the first unsynced account
	if iter.First() {
		// Ensure the key has the correct prefix
		if bytes.HasPrefix(iter.Key(), PrefixUnsyncedAccount) {
			publicKey := types.PublicKey(iter.Key()[1:])
			return &publicKey, nil
		}
	}

	return nil, errors.New("all accounts are synced")
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
	return utils.Serialize(a, binary.LittleEndian)
}

func (a *AccountRecord) Deserialize(data []byte) *AccountRecord {
	utils.Deserialize(data, a, binary.LittleEndian)
	return a
}

func SetAccount(publicKey [32]byte, account AccountRecord) {
	key := append(PrefixAccount, publicKey[:]...)
	if err := db.Set(key, account.Serialize(), pebble.Sync); err != nil {
		panic(err)
	}
}

func GetAccount(publicKey types.PublicKey) *AccountRecord {
	key := append(PrefixAccount, publicKey[:]...)
	serialized, closer, err := db.Get(key)

	if err != nil {
		return nil
	}

	if err := closer.Close(); err != nil {
		panic(err)
	}

	return (&AccountRecord{}).Deserialize(serialized)
}

func PutBlock(blockHash types.Hash, record BlockRecord) {
	key := append(PrefixBlock, blockHash[:]...)
	if err := db.Set(key, record.Serialize(), pebble.Sync); err != nil {
		panic(err)
	}

	fmt.Println("saved block:", hex.EncodeToString(blockHash[:]))
}

func GetBlock(blockHash [32]byte) *BlockRecord {
	key := append(PrefixBlock, blockHash[:]...)
	serialized, closer, err := db.Get(key)

	if err != nil {
		return nil
	}

	if err := closer.Close(); err != nil {
		panic(err)
	}

	return (&BlockRecord{}).Deserialize(serialized)
}

func GetLastBlockHash() [32]byte {
	iter, err := db.NewIter(&pebble.IterOptions{})
	defer iter.Close()

	if err != nil {
		panic(err)
	}

	if iter.Last() {
		return [32]byte(iter.Key())
	} else {
		return [32]byte{}
	}
}

func CountBlocks() uint64 {
	iter, _ := db.NewIter(&pebble.IterOptions{LowerBound: PrefixBlock})
	defer iter.Close()
	count := uint64(0)
	for iter.SeekGE(PrefixBlock); iter.Valid() && bytes.HasPrefix(iter.Key(), PrefixBlock); iter.Next() {
		count++
	}
	return count
}
