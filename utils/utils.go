package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
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
