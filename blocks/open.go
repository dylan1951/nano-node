package blocks

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/blake2b"
)

type OpenBlock struct {
	Source         [32]byte
	Representative [32]byte
	Account        [32]byte
	Signature      [64]byte
	Work           uint64
}

func (b *OpenBlock) Hash() [32]byte {
	var buf bytes.Buffer
	buf.Write(b.Source[:])
	buf.Write(b.Representative[:])
	buf.Write(b.Account[:])
	return blake2b.Sum256(buf.Bytes())
}

func (b *OpenBlock) Print() {
	fmt.Printf("Source:         %s\n", hex.EncodeToString(b.Source[:]))
	fmt.Printf("Representative: %s\n", hex.EncodeToString(b.Representative[:]))
	fmt.Printf("Account:        %s\n", hex.EncodeToString(b.Account[:]))
	fmt.Printf("Signature:      %s\n", hex.EncodeToString(b.Signature[:]))
	fmt.Printf("Work:           %x\n", b.Work)
}
