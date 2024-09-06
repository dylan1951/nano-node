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

type OpenBlock struct {
	Source         types.Hash
	Representative types.PublicKey
	Account        types.PublicKey
	BlockCommon
}

func (b *OpenBlock) Hash() types.Hash {
	var buf bytes.Buffer
	buf.Write(b.Source[:])
	buf.Write(b.Representative[:])
	buf.Write(b.Account[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *OpenBlock) Read(r io.Reader) *OpenBlock {
	io.ReadFull(r, b.Source[:])
	io.ReadFull(r, b.Representative[:])
	io.ReadFull(r, b.Account[:])
	binary.Read(r, binary.LittleEndian, &b.BlockCommon)
	return b
}

func (b *OpenBlock) Print() {
	fmt.Printf("Source:         %s\n", hex.EncodeToString(b.Source[:]))
	fmt.Printf("Representative: %s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("AccountRecord:        %s\n", hex.EncodeToString(b.Account[:]))
	fmt.Printf("Signature:      %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:           %x\n", b.Work)
}

func (b *OpenBlock) Serialize() []byte {
	return append([]byte{byte(LegacyOpen)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *OpenBlock) Type() Type {
	return LegacyOpen
}

func (b *OpenBlock) GetPrevious() types.Hash {
	return types.Hash{}
}
