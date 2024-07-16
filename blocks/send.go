package blocks

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"node/utils"
)

type SendBlock struct {
	Previous    [32]byte
	Destination [32]byte
	Balance     [16]byte
	Signature   [64]byte
	Work        uint64
}

func (b *SendBlock) Print() {
	fmt.Printf("Previous:    %s\n", hex.EncodeToString(b.Previous[:]))
	fmt.Printf("Destination: %s\n", hex.EncodeToString(b.Destination[:]))
	fmt.Printf("Balance:     %s\n", hex.EncodeToString(b.Balance[:]))
	fmt.Printf("Signature:   %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:        %x\n", b.Work)
}

func (b *SendBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Previous[:])
	buf.Write(b.Destination[:])
	buf.Write(b.Balance[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *SendBlock) Serialize() []byte {
	return append([]byte{byte(Send)}, utils.Serialize(b, binary.LittleEndian)...)
}

func (b *SendBlock) Type() Type {
	return Send
}
