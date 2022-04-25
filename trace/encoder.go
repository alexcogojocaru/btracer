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

func (e *Encoder) Compute() []byte {
	token := make([]byte, 16)
	rand.Seed(GetCurrentTimestamp())
	rand.Read(token)

	return token
}
