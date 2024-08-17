package messages

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"log"
	"node/utils"
)

const (
	Invalid = iota
	Blocks
	AccountInfo
	Frontiers
)

type AscPullReqHeader struct {
	Type byte
	Id   uint64
}

type AscPullReq struct{}

func ReadAscPullReq(reader io.Reader, extensions Extensions) *AscPullReq {
	header := &AscPullReqHeader{}
	if err := binary.Read(reader, binary.BigEndian, header); err != nil {
		log.Fatalf("Error reading AscPullReq header: %v", err)
	}

	switch header.Type {
	case Blocks:
		payload := &BlocksPayload{}
		if err := binary.Read(reader, binary.BigEndian, payload); err != nil {
			log.Fatalf("Error reading Blocks payload: %v", err)
		}
		//fmt.Printf("Start: %v\n", payload.Start)
		//fmt.Printf("Count: %d\n", payload.Count)
		//fmt.Printf("StartType: %v\n", payload.StartType)
	case AccountInfo:
	case Frontiers:
	default:
		log.Fatalf("Unknown AscPullReq type: %x", header.Type)
	}

	return &AscPullReq{}
}

type HashType byte

const (
	Account HashType = 0
	Block   HashType = 1
)

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
	msg := NewHeader(MsgAscPullReq, 34).Serialize()
	id := utils.Read[uint64](rand.Reader, binary.BigEndian)

	header := AscPullReqHeader{
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
	msg := NewHeader(MsgAscPullReq, 34).Serialize()
	id := utils.Read[uint64](rand.Reader, binary.LittleEndian)

	header := AscPullReqHeader{
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
