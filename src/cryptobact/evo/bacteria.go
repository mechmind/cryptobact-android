package evo

import "fmt"

type Bacteria struct {
    Chromosome *Chromosome
    TTL int
    Energy float64
    X float64
    Y float64
    Angle float64
    Born bool
}

func (b *Bacteria) String() string {
    return fmt.Sprintf("E%6.1f TTL%5d :: %s %t",
        b.Energy,
        b.TTL,
        b.Chromosome.DNA,
        b.Born,
    )
}
