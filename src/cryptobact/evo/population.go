package evo

import "log"
import "fmt"
import "math/rand"

var _ = fmt.Print
var _ = log.Print

type Population struct {
	// @TODO author
	Chain   *Chromochain
	Options *PopulationOptions

	Miner *Miner

	bacts []*Bacteria
}

type PopulationOptions struct {
	Attitudes           map[string]*Attitude
	MutateProbability   float64
	MutateRate          float64
	RecombinationChance float64
	RecombinationDrop   float64
}

var DefaultPopulationOptions = &PopulationOptions{
	Attitudes:           nil,
	MutateProbability:   0.5,
	MutateRate:          1.0,
	RecombinationChance: 1.0,
	RecombinationDrop:   10,
}

type Attitude struct {
	Pattern string
	Max     uint
}

func NewPopulation(miner *Miner, chain *Chromochain,
	options *PopulationOptions) *Population {
	if options == nil {
		options = DefaultPopulationOptions
	}

	bacts := make([]*Bacteria, 0)
	chromos := chain.GetChromosomes()
	for _, c := range chromos {
		bacts = append(bacts, &Bacteria{Chromosome: c, Born: true})
	}

	return &Population{bacts: bacts, Options: options, Miner: miner,
		Chain: chain}
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

	new_bacteria := &Bacteria{
		Chromosome: &Chromosome{Author: a.Chromosome.Author, DNA: new_dna},
		Born:       false}

	// mining here!
	p.Miner.Mine(new_bacteria.Chromosome)

	p.bacts = append(p.bacts, new_bacteria)

	return new_bacteria
}

func (p *Population) GetBacts() []*Bacteria {
	return p.bacts
}

func (p *Population) GetAttitude(b *Bacteria, attitude_id string) uint {
	if p.Options.Attitudes == nil {
		return 0
	}

	return b.Chromosome.DNA.MatchPatternCount(
		p.Options.Attitudes[attitude_id].Pattern)
}

func (p *Population) GetGene(b *Bacteria, index uint) float64 {
	genes := b.Chromosome.DNA.Genes()
	return float64(genes[int(index)%len(genes)])
}

func (p *Population) Kill(target *Bacteria) {
	alive := make([]*Bacteria, 0)

	for _, b := range p.bacts {
		if b != target {
			alive = append(alive, b)
		}
	}

	p.bacts = alive
}

func (p *Population) CatchNewBorn() bool {
	newborn := p.Miner.Extract()
	if newborn != nil {
		for _, b := range p.bacts {
			if b.Chromosome == newborn {
				b.Born = true
				return true
			}
		}
	}

	return false
}
