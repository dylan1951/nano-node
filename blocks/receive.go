package blocks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"node/utils"
)

type ReceiveBlock struct {
	Previous  [32]byte
	Source    [32]byte
	Signature [64]byte
	Work      uint64
}

func (b *ReceiveBlock) Print() {
	fmt.Printf("Previous:  %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Source:    %s\n", hex.EncodeToString(b.Source[:]))
	fmt.Printf("Signature: %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work: 	  %x\n", b.Work)
}

func (b *ReceiveBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Source[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *ReceiveBlock) Serialize() []byte {
	return append([]byte{byte(Receive)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *ReceiveBlock) Type() Type {
	return Receive
}
