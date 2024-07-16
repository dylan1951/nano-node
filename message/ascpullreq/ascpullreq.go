package ascpullreq

import (
	"crypto/rand"
	"encoding/binary"
	"node/message"
	"node/utils"
)

type PullType byte

const (
	Invalid     PullType = 0x0
	Blocks      PullType = 0x1
	AccountInfo PullType = 0x2
	Frontiers   PullType = 0x3
)

type HashType byte

const (
	Account = 0
	Block   = 1
)

type Header struct {
	Type PullType
	Id   uint64
}

type BlocksPayload struct {
	Start     [32]byte
	Count     uint8
	StartType HashType
}

type FrontiersPayload struct {
	Start [32]byte
	Count uint16
}

func BlocksRequest(start [32]byte, startType HashType) []byte {
	msg := message.NewHeader(message.AscPullReq, 34).Serialize()
	id := utils.Read[uint64](rand.Reader, binary.BigEndian)

	header := Header{
		Type: Blocks,
		Id:   *id,
	}

	payload := BlocksPayload{
		Start:     start,
		Count:     1,
		StartType: startType,
	}

	msg = append(msg, utils.Serialize(header, binary.BigEndian)...)
	msg = append(msg, utils.Serialize(payload, binary.BigEndian)...)

	return msg
}

func FrontiersRequest(start [32]byte, count uint16) []byte {
	msg := message.NewHeader(message.AscPullReq, 34).Serialize()
	id := utils.Read[uint64](rand.Reader, binary.LittleEndian)

	header := Header{
		Type: Frontiers,
		Id:   *id,
	}

	payload := FrontiersPayload{
		Start: start,
		Count: count,
	}

	msg = append(msg, utils.Serialize(header, binary.BigEndian)...)
	msg = append(msg, utils.Serialize(payload, binary.BigEndian)...)

	return msg
}
