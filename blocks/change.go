package blocks

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
)

type ChangeBlock struct {
	Previous       [32]byte
	Representative [32]byte
	Signature      [64]byte
	Work           uint64
}

func (b *ChangeBlock) Print() {
	//TODO implement me
	panic("implement me")
}

func (b *ChangeBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Representative[:])
	return blake2b.Sum256(buf.Bytes())
}
