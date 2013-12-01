package main

import "log"
//import "time"
import "cryptobact/engine"
//import "cryptobact/infektor"

type Updater struct {}

func (f Updater) Update(w *engine.World) {
	return
}

func main() {
    log.Println("testing engine")
    u := Updater{}
    engine.Loop(u)
}
