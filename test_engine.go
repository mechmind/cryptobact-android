package main

import "log"
import "./src/cryptobact/engine"

type Updater struct {}

func (f Updater) Update(w *engine.World) {
	return
}

func main() {
	log.Println("testing engine")
	u := Updater{}
	engine.Loop(u)
	log.Println("done")
}
