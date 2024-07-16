package blocks

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
)

type SendBlock struct {
	Previous    [32]byte
	Destination [32]byte
	Balance     [16]byte
	Signature   [64]byte
	Work        uint64
}

func (b *SendBlock) Print() {
	//TODO implement me
	panic("implement me")
}

func (b *SendBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Destination[:])
	buf.Write(b.Balance[:])
	return blake2b.Sum256(buf.Bytes())
}
