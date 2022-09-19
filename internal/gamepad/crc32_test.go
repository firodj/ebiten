package gamepad

import (
	"fmt"
	"hash/crc32"
	"testing"
)

func TestAja(t *testing.T) {
	fmt.Printf("0x%x\n", crc32.ChecksumIEEE([]byte("The quick brown fox jumps over the lazy dog")))
}