package main

import (
	"go.delic.rs/ligno"
)

func main() {
	ligno.Info("Some message", "key1", "value1", "key2", "value2")
	ligno.WaitAll()
}
