package main

import (
	"moussaud.org/pets/service/pets"

	. "moussaud.org/pets/internal"
)

func main() {
	LoadConfiguration()
	NewGlobalTracer()
	pets.Start()
}
