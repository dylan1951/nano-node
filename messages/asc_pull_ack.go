package messages

import (
	"encoding/binary"
	"io"
	"log"
	"node/blocks"
	"node/types"
	"node/utils"
)

type AscPullAck struct {
	Frontiers []*Frontier
	Blocks    []blocks.Block
}

type Frontier struct {
	Account types.PublicKey
	Hash    types.Hash
}

func (f *Frontier) IsZero() bool {
	return f.Account == [32]byte{} || f.Hash == [32]byte{}
}

func ReadAscPullAck(r io.Reader, extensions Extensions) AscPullAck {
	type Header struct {
		Type byte
		Id   uint64
	}

	header := utils.Read[Header](r, binary.BigEndian)
	ack := AscPullAck{}

	switch header.Type {
	case Blocks:
		ack.Blocks = make([]blocks.Block, 0)

		for b := blocks.Read(r); b != nil; b = blocks.Read(r) {
			ack.Blocks = append(ack.Blocks, b)
		}
	case Frontiers:
		ack.Frontiers = make([]*Frontier, 0)

		for f := utils.Read[Frontier](r, binary.BigEndian); !f.IsZero(); f = utils.Read[Frontier](r, binary.BigEndian) {
			ack.Frontiers = append(ack.Frontiers, f)
		}
	default:
		log.Fatalf("Unsupported AscPullAck PullType: %x", header.Type)
	}

	return ack
}
