package types

import (
	"encoding/hex"
	"encoding/json"
)

type Hash [32]byte
type Signature [64]byte
type PublicKey = Hash

func (h Hash) GoString() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) IsZero() bool {
	return h == Hash{}
}

func (h Hash) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(h.GoString())
	return jsonData, err
}

func (s Signature) GoString() string {
	return hex.EncodeToString(s[:])
}

func (s Signature) MarshalJSON() ([]byte, error) {
	jsonData, err := json.Marshal(s.GoString())
	return jsonData, err
}
