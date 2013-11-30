package population

import "math/rand"

import . "cryptobact/evo/chromochain"
import . "cryptobact/evo/chromosome"
import . "cryptobact/evo/bacteria"
import . "cryptobact/evo/dna"

type Population struct {
    //author Author
    Bacts []*Bacteria
}

func NewPopulation(chain *Chain) *Population {
    bacts := make([]*Bacteria, 0)
    chromos := chain.GetLastChromosomes()
    for _, c := range chromos {
        bacts = append(bacts, &Bacteria{Chromosome: c})
    }

    return &Population{Bacts: bacts}
}

func (p *Population) Fuck(a *Bacteria, b *Bacteria) *Bacteria {
    new_dna := Crossover(a.Chromosome.DNA, b.Chromosome.DNA)

    // @FIXME hardcode
    new_dna.Mutate(0.5, 1)

    user_choice := map[string]uint{
        "11.1": 6,
        "00.0": 2,
        "0101": 2,
    }

    second_recomb_change := 1.0
    for pattern, count := range user_choice {
        if new_dna.MatchPatternCount(pattern) >= count {
            continue
        }

        if rand.Float64() >= second_recomb_change {
            new_dna.Recombine(pattern)
            second_recomb_change /= 10
        }
    }

    // @TBD mining here
    new_bact := &Bacteria{Chromosome: &Chromosome{DNA: new_dna}}
    p.Bacts = append(p.Bacts, new_bact)

    return new_bact
}
