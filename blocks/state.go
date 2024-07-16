package blocks

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"node/utils"
)

type StateBlock struct {
	Account        [32]byte
	Previous       [32]byte
	Representative [32]byte
	Balance        [16]byte
	Link           [32]byte
	Signature      [64]byte
	Work           uint64
}

func (b *StateBlock) Hash() [32]byte {
	preamble := [32]byte{}
	preamble[31] = byte(State)
	hash, _ := blake2b.New256(nil)
	hash.Write(preamble[:])
	hash.Write(b.Account[:])
	hash.Write(b.Previous[:])
	hash.Write(b.Representative[:])
	hash.Write(b.Balance[:])
	hash.Write(b.Link[:])
	return [32]byte(hash.Sum(nil))
}

func (b *StateBlock) Print() {
	fmt.Printf("Account: 		%s\n", hex.EncodeToString(b.Account[:]))
	fmt.Printf("Previous: 		%s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Representative: 	%s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("Balance:        	%s\n", hex.EncodeToString(b.Balance[:]))
	fmt.Printf("Link:        	%s\n", hex.EncodeToString(b.Link[:]))
	fmt.Printf("Signature:       %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:           	%x\n", b.Work)
}

func (b *StateBlock) Serialize() []byte {
	return append([]byte{byte(State)}, utils.Serialize(b, binary.BigEndian)...)
}

func (b *StateBlock) Type() Type {
	return State
}
