package main

import (
	"moussaud.org/dogs/service/dogs"

	. "moussaud.org/dogs/internal"
)

func main() {
	LoadConfiguration()
	NewGlobalTracer()
	dogs.Start()
}
