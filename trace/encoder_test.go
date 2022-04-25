package trace_test

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/alexcogojocaru/btracer/trace"
)

func TestEncoderHashing(t *testing.T) {
	token := make([]byte, 16)
	rand.Seed(trace.GetCurrentTimestamp())
	rand.Read(token)

	t.Logf("bytes=%x", token)
	t.Logf("string=%s", hex.EncodeToString(token))
}
