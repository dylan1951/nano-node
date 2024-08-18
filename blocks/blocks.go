package blocks

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"node/types"
	"node/utils"
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
	BlockCommon() *BlockCommon
}

type BlockCommon struct {
	Signature [64]byte
	Work      uint64
}

func (bc *BlockCommon) BlockCommon() *BlockCommon {
	return bc
}

func Read(r io.Reader) Block {
	var blockType Type

	if err := binary.Read(r, binary.BigEndian, &blockType); err != nil {
		log.Fatalf("Failed to read block type: %v", err)
	}

	switch blockType {
	case LegacyOpen:
		return utils.Read[OpenBlock](r, binary.LittleEndian)
	case LegacySend:
		return utils.Read[SendBlock](r, binary.LittleEndian)
	case LegacyReceive:
		return utils.Read[ReceiveBlock](r, binary.LittleEndian)
	case LegacyChange:
		return utils.Read[ChangeBlock](r, binary.LittleEndian)
	case State:
		return utils.Read[StateBlock](r, binary.BigEndian)
	default:
		panic(fmt.Sprintf("Unknown block type: %d", blockType))
	}

	return nil
}

func Deserialize(serialized []byte) Block {
	return Read(bytes.NewReader(serialized))
}
