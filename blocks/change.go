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

type ChangeBlock struct {
	Previous       [32]byte
	Representative [32]byte
	BlockCommon
}

func (b *ChangeBlock) Print() {
	fmt.Printf("Previous:       %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Representative: %s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("Signature:      %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:           %x\n", b.Work)
}

func (b *ChangeBlock) Read(r io.Reader) *ChangeBlock {
	io.ReadFull(r, b.Previous[:])
	io.ReadFull(r, b.Representative[:])
	binary.Read(r, binary.LittleEndian, &b.BlockCommon)
	return b
}

func (b *ChangeBlock) Hash() types.Hash {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Representative[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *ChangeBlock) Serialize() []byte {
	return append([]byte{byte(LegacyChange)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *ChangeBlock) Type() Type {
	return LegacyChange
}

func (b *ChangeBlock) GetPrevious() types.Hash {
	return b.Previous
}
