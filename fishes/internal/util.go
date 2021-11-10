package internal

import (
	"fmt"
	"math/rand"
	"time"
)

var RAND = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewOrderNum() string {
	b := [16]byte{}
	RAND.Read(b[:])
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func RandSimDelay() {
	SimDelayChance := 0.3333
	SimDelayMS := 1000
	if RAND.Float64() < SimDelayChance {
		time.Sleep(time.Duration(RAND.Intn(SimDelayMS)) * time.Millisecond)
	}
}
