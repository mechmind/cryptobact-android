package bacteria

import . "cryptobact/evo/chromosome"

type Bacteria struct {
    Chromosome *Chromosome
    TTL uint
    Energy int
    X int
    Y int
}
