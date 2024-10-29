package main


import (
	"crypto/rand"
	"fmt"
	"io"
)

func GenerateUUID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	// Set version (4) and variant (2)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

