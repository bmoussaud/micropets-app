package main

import (
	"moussaud.org/cats/service/cats"

	. "moussaud.org/cats/internal"
)

func main() {
	LoadConfiguration()
	NewGlobalTracer()
	cats.Start()
}
