package messages

import (
	"encoding/binary"
	"io"
	"node/types"
	"node/utils"
)

type TelemetryAck struct {
	Signature types.Signature
	TelemetryData
}

type TelemetryData struct {
	NodeId            types.PublicKey
	BlockCount        uint64
	CementedCount     uint64
	UncheckedCount    uint64
	AccountCount      uint64
	BandwidthCap      uint64 // bits per second. 0 is unlimited
	PeerCount         uint32
	ProtocolVersion   uint8
	Uptime            uint64 // seconds since started
	GenesisBlock      types.Hash
	MajorVersion      byte
	MinorVersion      byte
	PatchVersion      byte
	PrereleaseVersion byte
	Maker             byte   // implementation identifier
	Timestamp         uint64 // current unix epoch in milliseconds
	ActiveDifficulty  uint64
}

func ReadTelemetryAck(r io.Reader, extensions Extensions) TelemetryAck {
	return *utils.Read[TelemetryAck](r, binary.BigEndian)
}
