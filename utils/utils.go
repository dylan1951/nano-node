package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
)

func ReverseBytes(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

func MustDecodeHex(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return bytes
}

func MustDecodeHex32(s string) [32]byte {
	var res [32]byte
	bytes := MustDecodeHex(s)
	copy(res[:], bytes)
	return res
}

func MustDecodeHex64(s string) [64]byte {
	var res [64]byte
	bytes := MustDecodeHex(s)
	copy(res[:], bytes)
	return res
}

func Serialize(obj interface{}, order binary.ByteOrder) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, order, obj); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func Deserialize(data []byte, obj interface{}, order binary.ByteOrder) {
	if err := binary.Read(bytes.NewReader(data), order, obj); err != nil {
		log.Fatal(err)
	}
}

func Read[T any](reader io.Reader, order binary.ByteOrder) *T {
	ret := new(T)
	if err := binary.Read(reader, order, ret); err != nil {
		log.Fatal(err)
	}
	return ret
}

func Write(writer io.Writer, value interface{}) {
	if err := binary.Write(writer, binary.BigEndian, value); err != nil {
		log.Fatal(err)
	}
}

func PrettyPrint(i interface{}) {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(s))
}

func byteArrayToPercentage(bytes [32]byte) float64 {
	// Convert the byte array to a big.Int
	bigIntValue := new(big.Int).SetBytes(bytes[:])

	// Create a big.Int representing the maximum value (2^256 - 1)
	maxValue := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)

	// Create big.Float versions of bigIntValue and maxValue
	bigFloatValue := new(big.Float).SetInt(bigIntValue)
	bigFloatMax := new(big.Float).SetInt(maxValue)

	// Calculate the percentage: (bigFloatValue / bigFloatMax) * 100
	percentage := new(big.Float).Quo(bigFloatValue, bigFloatMax)
	percentage.Mul(percentage, big.NewFloat(100))

	// Convert to float64
	result, _ := percentage.Float64()
	return result
}
