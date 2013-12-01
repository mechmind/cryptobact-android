package main

import "log"
//import "time"
import "cryptobact/engine"
//import "cryptobact/infektor"

type Updater struct {}

func (f Updater) Update(w *engine.World) {
    for _, p := range w.Populations {
        for _, b := range p.GetBacts() {
            log.Println(b)
        }
    }
}

func main() {
    log.Println("testing engine")
    u := Updater{}
    engine.Loop(u)
}
