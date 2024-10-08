package blocks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"io"
	"node/types"
	"node/utils"
)

type ReceiveBlock struct {
	Previous [32]byte
	Source   [32]byte
	BlockCommon
}

func (b *ReceiveBlock) Print() {
	fmt.Printf("Previous:  %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Source:    %s\n", hex.EncodeToString(b.Source[:]))
	fmt.Printf("Signature: %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work: 	  %x\n", b.Work)
}

func (b *ReceiveBlock) Read(r io.Reader) *ReceiveBlock {
	io.ReadFull(r, b.Previous[:])
	io.ReadFull(r, b.Source[:])
	binary.Read(r, binary.LittleEndian, &b.BlockCommon)
	return b
}

func (b *ReceiveBlock) Hash() types.Hash {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Source[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *ReceiveBlock) Serialize() []byte {
	return append([]byte{byte(LegacyReceive)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *ReceiveBlock) Type() Type {
	return LegacyReceive
}

func (b *ReceiveBlock) GetPrevious() types.Hash {
	return b.Previous
}
