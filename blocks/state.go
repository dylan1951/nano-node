package blocks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"io"
	"node/types"
	"node/types/uint128"
	"node/utils"
)

type StateBlock struct {
	Account        types.PublicKey
	Previous       types.Hash
	Representative types.PublicKey
	Balance        uint128.Uint128
	Link           types.Hash
	BlockCommon
}

func (b *StateBlock) Hash() types.Hash {
	preamble := [32]byte{31: byte(State)}
	var buf bytes.Buffer
	buf.Grow(176)
	buf.Write(preamble[:])
	buf.Write(b.Account[:])
	buf.Write(b.Previous[:])
	buf.Write(b.Representative[:])
	buf.Write(b.Balance.BytesBE())
	buf.Write(b.Link[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *StateBlock) Read(r io.Reader) *StateBlock {
	io.ReadFull(r, b.Account[:])
	io.ReadFull(r, b.Previous[:])
	io.ReadFull(r, b.Representative[:])
	b.Balance = uint128.ReadBE(r)
	io.ReadFull(r, b.Link[:])
	binary.Read(r, binary.BigEndian, &b.BlockCommon)
	return b
}

func (b *StateBlock) Print() {
	fmt.Printf("Hash: 		%s\n", b.Hash().GoString())
	fmt.Printf("Account: 		%s\n", hex.EncodeToString(b.Account[:]))
	fmt.Printf("Previous: 		%s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Representative: 	%s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("Balance:        	%s\n", b.Balance.String())
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

func (b *StateBlock) GetPrevious() types.Hash {
	return b.Previous
}
