package blocks

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
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

func (b *StateBlock) Print() {
	//TODO implement me
	panic("implement me")
}

func (b *StateBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Account[:])
	buf.Write(b.Previous[:])
	buf.Write(b.Representative[:])
	buf.Write(b.Balance[:])
	buf.Write(b.Link[:])
	return blake2b.Sum256(buf.Bytes())
}
