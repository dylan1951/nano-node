package asc_pull_req

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"node/messages"
	"node/utils"
)

type AscPullReq struct {
}

func Read(reader io.Reader, extensions uint16) *AscPullReq {
	header := &Header{}
	if err := binary.Read(reader, binary.BigEndian, header); err != nil {
		log.Fatalf("Error reading AscPullReq header: %v", err)
	}

	switch header.Type {
	case Blocks:
		payload := &BlocksPayload{}
		if err := binary.Read(reader, binary.BigEndian, payload); err != nil {
			log.Fatalf("Error reading Blocks payload: %v", err)
		}
		fmt.Printf("Start: %v\n", payload.Start)
		fmt.Printf("Count: %d\n", payload.Count)
		fmt.Printf("StartType: %v\n", payload.StartType)
	case AccountInfo:
	case Frontiers:
	default:
		log.Fatalf("Unknown AscPullReq type: %x", header.Type)
	}

	return &AscPullReq{}
}

type PullType byte

const (
	Invalid     PullType = 0x0
	Blocks      PullType = 0x1
	AccountInfo PullType = 0x2
	Frontiers   PullType = 0x3
)

type HashType byte

const (
	Account HashType = 0
	Block   HashType = 1
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
	msg := messages.NewHeader(messages.MsgAscPullReq, 34).Serialize()
	id := utils.Read[uint64](rand.Reader, binary.BigEndian)

	header := Header{
		Type: Blocks,
		Id:   *id,
	}

	payload := BlocksPayload{
		Start:     start,
		Count:     50,
		StartType: startType,
	}

	msg = append(msg, utils.Serialize(header, binary.BigEndian)...)
	msg = append(msg, utils.Serialize(payload, binary.BigEndian)...)

	return msg
}

func FrontiersRequest(start [32]byte, count uint16) []byte {
	msg := messages.NewHeader(messages.MsgAscPullReq, 34).Serialize()
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
