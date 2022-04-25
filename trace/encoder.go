package trace

import (
	"hash"
	"time"
)

type Encoder struct {
	Hash hash.Hash
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func (e *Encoder) Compute(item ...string) []byte {
	encodedString := ""
	for _, element := range item {
		encodedString += element
	}

	e.Hash.Sum([]byte(encodedString))
	encoded := e.Hash.Sum(nil)
	e.Hash.Reset()

	return encoded
}
