package main

import (
	"moussaud.org/fishes/service/fishes"

	. "moussaud.org/fishes/internal"
)

func main() {
	LoadConfiguration()
	NewGlobalTracer()
	fishes.Start()
}
