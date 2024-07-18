package telemetry_ack

import (
	"encoding/binary"
	"io"
	"node/types"
	"node/utils"
)

type TelemetryAck struct {
	Signature         types.Signature
	NodeId            types.PublicKey
	BlockCount        uint64
	CementedCount     uint64
	UncheckedCount    uint64
	AccountCount      uint64
	BandwidthCap      uint64
	PeerCount         uint32
	ProtocolVersion   uint8
	Uptime            uint64
	GenesisBlock      types.Hash
	MajorVersion      byte
	MinorVersion      byte
	PatchVersion      byte
	PrereleaseVersion byte
	Maker             byte
	Timestamp         uint64
	ActiveDifficulty  uint64
}

func Read(reader io.Reader, extensions uint16) *TelemetryAck {
	return utils.Read[TelemetryAck](reader, binary.BigEndian)
}
