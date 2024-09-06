package blocks

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"node/types"
)

type Type uint8

const (
	NotABlock     Type = 1
	LegacySend    Type = 2
	LegacyReceive Type = 3
	LegacyOpen    Type = 4
	LegacyChange  Type = 5
	State         Type = 6
)

type Block interface {
	Print()
	Hash() types.Hash
	Serialize() []byte
	Type() Type
	Common() *BlockCommon
	GetPrevious() types.Hash
}

type BlockCommon struct {
	Signature [64]byte
	Work      uint64
}

func (bc *BlockCommon) Common() *BlockCommon {
	return bc
}

func Read(r io.Reader) Block {
	var blockType Type

	if err := binary.Read(r, binary.BigEndian, &blockType); err != nil {
		log.Fatalf("Failed to read block type: %v", err)
	}

	switch blockType {
	case LegacyOpen:
		return (&OpenBlock{}).Read(r)
	case LegacySend:
		return (&SendBlock{}).Read(r)
	case LegacyReceive:
		return (&ReceiveBlock{}).Read(r)
	case LegacyChange:
		return (&ChangeBlock{}).Read(r)
	case State:
		return (&StateBlock{}).Read(r)
	case NotABlock:
		return nil
	default:
		panic(fmt.Sprintf("Unknown block type: %d", blockType))
	}
}

func Deserialize(serialized []byte) Block {
	return Read(bytes.NewReader(serialized))
}
