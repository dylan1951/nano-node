package blocks

import (
	"encoding/binary"
	"log"
	"net"
	"node/types"
	"node/utils"
)

type BlockType byte

const (
	NotABlock BlockType = 1
	Send      BlockType = 2
	Receive   BlockType = 3
	Open      BlockType = 4
	Change    BlockType = 5
	State     BlockType = 6
)

//goland:noinspection SpellCheckingInspection
var TestGenesisBlock = OpenBlock{
	Source:         utils.MustDecodeHex32("45C6FF9D1706D61F0821327752671BDA9F9ED2DA40326B01935AB566FB9E08ED"),
	Representative: types.MustParseAddress("nano_1jg8zygjg3pp5w644emqcbmjqpnzmubfni3kfe1s8pooeuxsw49fdq1mco9j"),
	Account:        types.MustParseAddress("nano_1jg8zygjg3pp5w644emqcbmjqpnzmubfni3kfe1s8pooeuxsw49fdq1mco9j"),
	Signature:      utils.MustDecodeHex64("15049467CAEE3EC768639E8E35792399B6078DA763DA4EBA8ECAD33B0EDC4AF2E7403893A5A602EB89B978DABEF1D6606BB00F3C0EE11449232B143B6E07170E"),
	Work:           0xbc1ef279c1a34eb1,
}

//goland:noinspection SpellCheckingInspection
var BetaGenesisBlock = OpenBlock{
	Source:         utils.MustDecodeHex32("259A438A8F9F9226130C84D902C237AF3E57C0981C7D709C288046B110D8C8AC"),
	Representative: types.MustParseAddress("nano_1betag7az9wk6rbis38s1d35hdsycz1bi95xg4g4j148p6afjk7embcurda4"),
	Account:        types.MustParseAddress("nano_1betag7az9wk6rbis38s1d35hdsycz1bi95xg4g4j148p6afjk7embcurda4"),
	Signature:      utils.MustDecodeHex64("BC588273AC689726D129D3137653FB319B6EE6DB178F97421D11D075B46FD52B6748223C8FF4179399D35CB1A8DF36F759325BD2D3D4504904321FAFB71D7602"),
	Work:           0xe87a3ce39b43b84c,
}

type Block interface {
	Print()
	Hash() [32]byte
}

func Read(conn net.Conn) Block {
	var blockType BlockType

	if err := binary.Read(conn, binary.BigEndian, &blockType); err != nil {
		log.Fatalf("Failed to read block type: %v", err)
	}

	var block Block

	switch blockType {
	case Open:
		block = utils.Read[OpenBlock](conn, binary.LittleEndian)
	case Send:
		block = utils.Read[SendBlock](conn, binary.LittleEndian)
	case Receive:
		block = utils.Read[ReceiveBlock](conn, binary.LittleEndian)
	case Change:
		block = utils.Read[ChangeBlock](conn, binary.LittleEndian)
	case State:
		block = utils.Read[StateBlock](conn, binary.BigEndian)
	case NotABlock:
		return nil
	default:
		log.Fatalf("Unknown block type: %v", blockType)
	}

	return block
}
