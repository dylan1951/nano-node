package main

import (
	"encoding/hex"
	"github.com/accept-nano/ed25519-blake2b"
	"github.com/joho/godotenv"
	"log"
	"node/config"
	"node/node"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	seed, err := hex.DecodeString(os.Getenv("SEED"))
	if err != nil || len(seed) != ed25519.SeedSize {
		log.Fatalf("Invalid SEED in .env file")
	}

	network := os.Getenv("NETWORK")

	n := node.NewNode(config.Config{
		PrivateKey: ed25519.NewKeyFromSeed(seed),
		PublicKey:  ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey),
		Network:    config.Networks[network],
	})

	n.Bootstrap()
	n.Listen()
}
