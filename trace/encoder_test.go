package trace_test

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/alexcogojocaru/btracer/trace"
)

func TestEncoderHashing(t *testing.T) {
	encoder := trace.Encoder{
		Hash: sha1.New(),
	}
	token := encoder.Compute()

	t.Logf("bytes=%x", token)
	t.Logf("string=%s", hex.EncodeToString(token))
}
