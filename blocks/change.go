package blocks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"node/utils"
)

type ChangeBlock struct {
	Previous       [32]byte
	Representative [32]byte
	Signature      [64]byte
	Work           uint64
}

func (b *ChangeBlock) Print() {
	fmt.Printf("Previous:       %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Representative: %s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("Signature:      %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:           %x\n", b.Work)
}

func (b *ChangeBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Representative[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *ChangeBlock) Serialize() []byte {
	return append([]byte{byte(Change)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *ChangeBlock) Type() Type {
	return Change
}
