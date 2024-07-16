package blocks

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
)

type ReceiveBlock struct {
	Previous  [32]byte
	Source    [32]byte
	Signature [64]byte
	Work      uint64
}

func (b *ReceiveBlock) Print() {
	//TODO implement me
	panic("implement me")
}

func (b *ReceiveBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Source[:])
	return blake2b.Sum256(buf.Bytes())
}
