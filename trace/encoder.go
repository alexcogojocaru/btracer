package trace

import (
	"hash"
	"math/rand"
	"time"
)

type Encoder struct {
	Hash hash.Hash
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func (e *Encoder) Compute(size int8) []byte {
	token := make([]byte, size)
	rand.Seed(GetCurrentTimestamp())
	rand.Read(token)

	return token
}
