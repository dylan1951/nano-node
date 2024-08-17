package config

import (
	"encoding/hex"
	"github.com/accept-nano/ed25519-blake2b"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var PrivateKey ed25519.PrivateKey
var PublicKey ed25519.PublicKey
var Network NetworkDetail

var ProtocolVersionMax = uint8(21)
var ProtocolVersionUsing = uint8(21)
var ProtocolVersionMin = uint8(20)

type NetworkDetail struct {
	Id      byte
	Address string
	Port    uint16
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
	},
	"beta": {
		Id:      'B',
		Address: "peering-beta.nano.org",
		Port:    54000,
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
