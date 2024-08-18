package config

import (
	"encoding/hex"
	"github.com/accept-nano/ed25519-blake2b"
	"github.com/joho/godotenv"
	"log"
	"node/blocks"
	"node/types"
	"node/utils"
	"os"
)

var PrivateKey ed25519.PrivateKey
var PublicKey ed25519.PublicKey
var Network NetworkDetail

var ProtocolVersionMax = uint8(21)
var ProtocolVersionUsing = uint8(21)
var ProtocolVersionMin = uint8(20)

var ActiveDifficulty = uint64(0xFFFFF00000000000)

type NetworkDetail struct {
	Id      byte
	Address string
	Port    uint16
	Genesis blocks.OpenBlock
}

var Networks = map[string]NetworkDetail{
	"main": {
		Id:      'C',
		Address: "peering.nano.org",
		Port:    7075,
	},
	"test": {
		Id:      'X',
		Address: "peering-test.nano.org",
		Port:    17075,
		Genesis: blocks.OpenBlock{
			Source:         utils.MustDecodeHex32("45C6FF9D1706D61F0821327752671BDA9F9ED2DA40326B01935AB566FB9E08ED"),
			Representative: types.MustParseAddress("nano_1jg8zygjg3pp5w644emqcbmjqpnzmubfni3kfe1s8pooeuxsw49fdq1mco9j"),
			Account:        types.MustParseAddress("nano_1jg8zygjg3pp5w644emqcbmjqpnzmubfni3kfe1s8pooeuxsw49fdq1mco9j"),
			BlockCommon: blocks.BlockCommon{
				Signature: utils.MustDecodeHex64("15049467CAEE3EC768639E8E35792399B6078DA763DA4EBA8ECAD33B0EDC4AF2E7403893A5A602EB89B978DABEF1D6606BB00F3C0EE11449232B143B6E07170E"),
				Work:      0xbc1ef279c1a34eb1,
			},
		},
	},
	"beta": {
		Id:      'B',
		Address: "peering-beta.nano.org",
		Port:    54000,
		Genesis: blocks.OpenBlock{
			Source:         utils.MustDecodeHex32("259A438A8F9F9226130C84D902C237AF3E57C0981C7D709C288046B110D8C8AC"),
			Representative: types.MustParseAddress("nano_1betag7az9wk6rbis38s1d35hdsycz1bi95xg4g4j148p6afjk7embcurda4"),
			Account:        types.MustParseAddress("nano_1betag7az9wk6rbis38s1d35hdsycz1bi95xg4g4j148p6afjk7embcurda4"),
			BlockCommon: blocks.BlockCommon{
				Signature: utils.MustDecodeHex64("BC588273AC689726D129D3137653FB319B6EE6DB178F97421D11D075B46FD52B6748223C8FF4179399D35CB1A8DF36F759325BD2D3D4504904321FAFB71D7602"),
				Work:      0xe87a3ce39b43b84c,
			},
		},
	},
}

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	seed, err := hex.DecodeString(os.Getenv("SEED"))
	if err != nil || len(seed) != ed25519.SeedSize {
		log.Fatalf("Invalid SEED in .env file")
	}

	network := os.Getenv("NETWORK")

	PrivateKey = ed25519.NewKeyFromSeed(seed)
	PublicKey = ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	Network = Networks[network]
}
