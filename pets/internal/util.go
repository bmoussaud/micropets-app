package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type ErrorStatus struct {
	Error string `json:"error"`
}

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

func WriteError(w http.ResponseWriter, err string, statusCode int) []byte {
	log.Println(err)
	bytes, _ := json.Marshal(ErrorStatus{Error: err})
	w.WriteHeader(statusCode)
	w.Write(bytes)
	return bytes
}
