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

type SendBlock struct {
	Previous    types.Hash
	Destination types.PublicKey
	Balance     uint128.Uint128
	BlockCommon
}

func (b *SendBlock) Hash() types.Hash {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Destination[:])
	buf.Write(b.Balance.Bytes())
	return blake2b.Sum256(buf.Bytes())
}

func (b *SendBlock) Read(r io.Reader) *SendBlock {
	io.ReadFull(r, b.Previous[:])
	io.ReadFull(r, b.Destination[:])
	b.Balance = uint128.Read(r)
	binary.Read(r, binary.LittleEndian, &b.BlockCommon)
	return b
}

func (b *SendBlock) Print() {
	fmt.Printf("Previous:    %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Destination: %s\n", hex.EncodeToString(b.Destination[:]))
	fmt.Printf("Balance:     %s\n", b.Balance.String())
	fmt.Printf("Signature:   %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:        %x\n", b.Work)
}

func (b *SendBlock) Serialize() []byte {
	return append([]byte{byte(LegacySend)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *SendBlock) Type() Type {
	return LegacySend
}

func (b *SendBlock) GetPrevious() types.Hash {
	return b.Previous
}
