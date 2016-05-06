package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

func Random(size int) []byte {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	if err != nil {
	}
	return bytes
}

func Hex(data []byte) string {
	return hex.EncodeToString(data)
}

func RandomHex(size int) string {
	return Hex(Random(size))
}

func RandomNumber(min, max int) uint {
	p, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		fmt.Println(err.Error())
		return uint(min)
	}
	if uint(p.Uint64()) < uint(min) {
		return uint(min)
	} else {
		return uint(p.Uint64())
	}
}
