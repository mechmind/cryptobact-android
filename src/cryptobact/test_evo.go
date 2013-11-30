package main

import "fmt"
import "math/rand"
import . "cryptobact/evo/dna"
import . "cryptobact/evo/chromochain"
import . "cryptobact/evo/population"

func main() {
    chain := &Chain{}
    population := NewPopulation(chain)
    population.Fuck(population.Bacts[8], population.Bacts[9])

    for i, bact := range population.Bacts {
        fmt.Printf("%2d: %s\n", i + 1, bact.Chromosome.DNA)
    }
}

func test_dna() {
    dna_a := NewRandDNA(64)
    dna_b := NewRandDNA(64)

    fmt.Println("DNA A", dna_a)
    fmt.Println("DNA B", dna_b)

    dna_c := Crossover(dna_a, dna_b)

    fmt.Println("CROSSOVER")
    fmt.Println("DNA C", dna_c)

    dna_c.Mutate(0.5, 1)
    fmt.Println("MUT C", dna_c)
    fmt.Println("GENES", dna_c.Genes())

    patterns := []string{"111.1", "000.0"}
    for i := 0; i < 1000; i++ {
        for _, p := range patterns {
            if rand.Intn(10) != 1 {
                continue;
            }
            dna_c.Recombine(p)
            fmt.Println("REC C", dna_c, p, dna_c.MatchPatternCount(p))
            fmt.Println("GENES", dna_c.Genes())
        }
    }
}
