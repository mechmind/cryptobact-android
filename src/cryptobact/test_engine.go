package main

import(
	"cryptobact/engine"
	"log"
)

func main() {
	m := engine.GetMap()
	for i := range m {
		log.Println(m[i])
	}
}
