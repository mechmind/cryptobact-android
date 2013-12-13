package main

import (
	"cryptobact/engine"
	"log"
)

type Updater struct{}

func (f Updater) Update(w *engine.World) {
	//for _, p := range w.Populations {
	//    for _, b := range p.Bacts {
	//        //log.Println(b)
	//    }
	//}
}

func main() {
	log.Println("testing engine")
	u := Updater{}
	engine.Loop(u)
}
