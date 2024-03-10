package quantumid

import "testing"

func TestNewString(t *testing.T) {
	t.Logf("UUID style:       %v", UUID())
	t.Logf("UUID style(tidy): %v", UUIDTidy())
	t.Logf("Base58 style:     %v", Base58())
}
