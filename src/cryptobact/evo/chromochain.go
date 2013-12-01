package evo

type Chromochain struct {
    Author uint64
    Initial *Chromosome
    // @TBD
}

func (c *Chromochain) GetChromosomes() []*Chromosome {
    // @FIXME hardcode
    chromos := make([]*Chromosome, 0)
    var initial *DNA
    if c.Initial == nil {
        initial = NewRandDNA(64)
    } else {
        initial = c.Initial.DNA
    }
    for i := 0; i < 10; i++ {
        chromos = append(chromos, &Chromosome{DNA: initial,
            Author: c.Author})
    }

    return chromos
}
