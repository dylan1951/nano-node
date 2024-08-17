package blocks

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"node/utils"
)

type Type uint8

const (
	NotABlock Type = 1
	Send      Type = 2
	Receive   Type = 3
	Open      Type = 4
	Change    Type = 5
	State     Type = 6
)

type Block interface {
	Print()
	Hash() [32]byte
	Serialize() []byte
	Type() Type
}

func Read(r io.Reader) Block {
	var blockType Type

	if err := binary.Read(r, binary.BigEndian, &blockType); err != nil {
		log.Fatalf("Failed to read block type: %v", err)
	}

	switch blockType {
	case Open:
		return utils.Read[OpenBlock](r, binary.LittleEndian)
	case Send:
		return utils.Read[SendBlock](r, binary.LittleEndian)
	case Receive:
		return utils.Read[ReceiveBlock](r, binary.LittleEndian)
	case Change:
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
