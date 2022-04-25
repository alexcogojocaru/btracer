package trace

import (
	"fmt"
	"time"
)

type ReadTrace interface {
	GenerateID() []byte
}

type Trace struct {
}

func (t *Trace) GenerateID() []byte {
	timestamp := time.Now().Unix()
	encoded := []byte(fmt.Sprint(timestamp))
	return encoded
}
