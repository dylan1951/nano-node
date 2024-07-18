package asc_pull_ack

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"node/blocks"
	"node/types"
	"node/utils"
)

type PullType byte

const (
	Invalid     PullType = 0x0
	Blocks      PullType = 0x1
	AccountInfo PullType = 0x2
	Frontiers   PullType = 0x3
)

type Header struct {
	Type PullType
	Id   uint64
}

type AscPullAck struct {
	Frontiers []*Frontier
	Blocks    []*blocks.Block
}

type Frontier struct {
	Account types.PublicKey
	Hash    types.Hash
}

func (f *Frontier) IsZero() bool {
	return f.Account == [32]byte{} || f.Hash == [32]byte{}
}

func Read(reader io.Reader, extensions uint16) *AscPullAck {
	header := utils.Read[Header](reader, binary.BigEndian)
	ack := AscPullAck{}

	switch header.Type {
	case Blocks:
		ack.Blocks = make([]*blocks.Block, 1024)

		for b := blocks.Read(reader); b != nil; b = blocks.Read(reader) {
			hash := b.Hash()
			fmt.Println(hex.EncodeToString(hash[:]))
			ack.Blocks = append(ack.Blocks, &b)
		}
	case Frontiers:
		ack.Frontiers = make([]*Frontier, 1024)

		for f := utils.Read[Frontier](reader, binary.BigEndian); !f.IsZero(); f = utils.Read[Frontier](reader, binary.BigEndian) {
			fmt.Printf("frontier for %s is %s\n", f.Account.GoString(), f.Hash.GoString())
			ack.Frontiers = append(ack.Frontiers, f)
		}
	default:
		log.Fatalf("Unsupported AscPullAck PullType: %x", header.Type)
	}

	return &ack
}
