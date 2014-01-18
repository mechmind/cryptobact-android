package main

import (
	"fmt"
	"cryptobact/evo"
)

type TestCreature struct {
	C *evo.Chromosome
}

func main() {
	dnas := make([]*evo.DNA, 0)

	for i := 0; i < 2; i++ {
		dnas = append(dnas, evo.NewRandDNA(16, 1))
	    fmt.Println(dnas[i])
	}

	//for _, dna := range dnas {
	//    fmt.Println(dna)
	//}

	//for i := 0; i < 100; i++ {
	//    fmt.Println(dnas[0])

	//    dnas[0].Mutate(0.001, 10)
	//}
	//
	
	new_dna := evo.Crossover(dnas[0], dnas[1])

	fmt.Println(new_dna)
}
