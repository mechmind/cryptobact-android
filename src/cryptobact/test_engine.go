package main

import(
	"cryptobact/engine"
	"log"
)

func main() {
	log.Println("testing engine")
	engine.Loop()
	log.Println("done")
}
