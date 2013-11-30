package chromochain

import . "cryptobact/evo/chromosome"
import . "cryptobact/evo/dna"


type Chain struct {
    // @TBD
}

func (c *Chain) GetLastChromosomes() []*Chromosome {
    // @FIXME hardcode
    chromos := make([]*Chromosome, 0)
    for i := 0; i < 10; i++ {
        chromos = append(chromos, &Chromosome{DNA: NewRandDNA(64)})
    }

    return chromos
}
