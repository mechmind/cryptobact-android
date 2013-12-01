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
	Speed float64
	RotationSpeed float64
}

func (b *Bacteria) String() string {
<<<<<<< HEAD
    return fmt.Sprintf("E%6.1f TTL%5d :: %s %t",
=======
    return fmt.Sprintf("{%5.2f; %5.2f} A%3.2f E%6.1f TTL%5d :: %s [%t]\n",
        b.X, b.Y,
		b.Angle,
>>>>>>> 67dd007... Rotation implemented.
        b.Energy,
        b.TTL,
        b.Chromosome.DNA,
        b.Born,
    )
}
