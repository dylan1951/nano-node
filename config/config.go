package config

import "github.com/accept-nano/ed25519-blake2b"

type Config struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Network    Network
}

type Network struct {
	Id      byte
	Address string
	Port    uint16
}

var Networks = map[string]Network{
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
