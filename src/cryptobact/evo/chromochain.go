package evo

type Chromochain struct {
    // @TBD
}

func (c *Chromochain) GetLastChromosomes() []*Chromosome {
    // @FIXME hardcode
    chromos := make([]*Chromosome, 0)
    for i := 0; i < 10; i++ {
        chromos = append(chromos, &Chromosome{DNA: NewRandDNA(64)})
    }

    return chromos
}
