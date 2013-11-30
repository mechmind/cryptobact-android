package evo

import "fmt"

type Bacteria struct {
    Chromosome *Chromosome
    TTL uint
    Energy float64
    X float64
    Y float64
    Angle float64
    Born bool
}

func (b *Bacteria) String() string {
    return fmt.Sprintf("%p {%3.2f; %3.2f} E%3.2f TTL%d :: %s [%t]\n",
        b,
        b.X, b.Y,
        b.Energy,
        b.TTL,
        b.Chromosome.DNA,
        b.Born,
    )
}
