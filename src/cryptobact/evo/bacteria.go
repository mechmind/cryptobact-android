package evo

type Bacteria struct {
    Chromosome *Chromosome
    TTL uint
    Energy float64
    X float64
    Y float64
    Direction float64
    Born bool
}
