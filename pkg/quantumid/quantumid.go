package quantumid

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

func quantumID() []byte {
	raw := make([]byte, 16)
	s := time.Now().UnixNano()
	for i := 0; i < 8; i++ {
		raw[7-i] = byte(s % 256)
		s >>= 8
	}
	_, _ = rand.Read(raw[8:16])
	return raw
}

// Base58 new a quantum id in base58 style
func Base58() string {
	return base58.Encode(quantumID())
}

// UUID new a quantum id in UUID style
func UUID() string {
	raw := quantumID()
	return fmt.Sprintf("%x-%x-%x-%x-%x", raw[0:4], raw[4:6], raw[6:8], raw[8:10], raw[10:16])
}

// UUIDTidy new a quantum id in UUID style without -
func UUIDTidy() string {
	return strings.ReplaceAll(UUID(), "-", "")
}
