package evo

import "math/rand"

type Population struct {
    // @TODO author
    Options *PopulationOptions

    bacts []*Bacteria
}

type PopulationOptions struct {
    Attitudes map[string]*Attitude
    MutateProbability float64
    MutateRate float64
    RecombinationChance float64
    RecombinationDrop float64
}

var DefaultPopulationOptions = &PopulationOptions{
    Attitudes: nil,
    MutateProbability: 0.5,
    MutateRate: 1.0,
    RecombinationChance: 1.0,
    RecombinationDrop: 10,
}

type Attitude struct {
    Pattern string
    Max uint
}

func NewPopulation(chain *Chromochain, options *PopulationOptions) *Population {
    if options == nil {
        options = DefaultPopulationOptions
    }

    bacts := make([]*Bacteria, 0)
    chromos := chain.GetLastChromosomes()
    for _, c := range chromos {
        bacts = append(bacts, &Bacteria{Chromosome: c})
    }

    return &Population{bacts: bacts, Options: options}
}

func (p *Population) Fuck(a *Bacteria, b *Bacteria) *Bacteria {
    new_dna := Crossover(a.Chromosome.DNA, b.Chromosome.DNA)

    // @FIXME hardcode
    new_dna.Mutate(p.Options.MutateProbability, p.Options.MutateRate)

    second_recomb_change := p.Options.RecombinationChance
    for _, attitude := range p.Options.Attitudes {
        if new_dna.MatchPatternCount(attitude.Pattern) >= attitude.Max {
            continue
        }

        if rand.Float64() >= second_recomb_change {
            new_dna.Recombine(attitude.Pattern)
            second_recomb_change /= p.Options.RecombinationDrop
        }
    }

    // @TBD mining here
    new_bacteria := &Bacteria{Chromosome: &Chromosome{DNA: new_dna}}

    p.bacts = append(p.bacts, new_bacteria)

    return new_bacteria
}

func (p *Population) GetBacts() []*Bacteria {
    return p.bacts
}

func (p *Population) GetAttitude(b *Bacteria, attitude_id string) uint {
    return b.Chromosome.DNA.MatchPatternCount(
        p.Options.Attitudes[attitude_id].Pattern)
}

func (p *Population) GetGenome(b *Bacteria) []uint {
    return b.Chromosome.DNA.Genes()
}

func (p *Population) Clean() []*Bacteria {
    corpses := make([]*Bacteria, 0)
    alive := make([]*Bacteria, 0)

    for _, b := range p.bacts {
        if b.Energy <= 0 {
            corpses = append(corpses, b)
        } else {
            alive = append(alive, b)
        }
    }

    return corpses
}
