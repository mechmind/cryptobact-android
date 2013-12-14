package evo

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

var _ = fmt.Print
var _ = log.Print

type TraitMap map[string]*Trait

type Population struct {
	// @TODO author
	Traits TraitMap
	Chain  *Chromochain
	Env    *Environment

	Bacts []*Bacteria
}

func (t TraitMap) String() string {
	result := make([]string, 0)

	for k, v := range t {
		result = append(result, fmt.Sprintf("[%s] %s", k, v))
	}

	return strings.Join(result, "\n")
}

type Trait struct {
	Pattern string
	Max     uint
}

func (t *Trait) String() string {
	return fmt.Sprintf("{%s}: --> %d", t.Pattern, t.Max)
}

func NewPopulation(chain *Chromochain, traits TraitMap,
	env *Environment) *Population {
	if env == nil {
		env = DefaultEnvironment
	}

	bacts := make([]*Bacteria, 0)
	chromos := chain.GetChromosomes()
	for _, c := range chromos {
		new_bacteria := NewBacteria(c)
		new_bacteria.Born = true
		bacts = append(bacts, new_bacteria)
	}

	return &Population{
		Traits: traits,
		Env:    env,
		Chain:  chain,
		Bacts:  bacts,
	}
}

func (p *Population) Slice(bacts []*Bacteria) *Population {
	return &Population{Bacts: bacts, Env: nil, Chain: p.Chain,
		Traits: p.Traits}
}

func (p *Population) GetGene(b *Bacteria, index int) float64 {
	return 1.0
}

func (p *Population) Fuck(a *Bacteria, b *Bacteria) *Bacteria {
	new_dna := Crossover(a.Chromosome.DNA, b.Chromosome.DNA)

	new_dna.Mutate(p.Env.MutateProbability, p.Env.MutateRate)

	second_recomb_change := p.Env.RecombinationChance
	for _, attitude := range p.Traits {
		if new_dna.MatchPatternCount(attitude.Pattern) >= attitude.Max {
			continue
		}

		if rand.Float64() >= second_recomb_change {
			new_dna.Recombine(attitude.Pattern)
			second_recomb_change /= p.Env.RecombinationDrop
		}
	}

	new_bacteria := NewBacteria(&Chromosome{Author: a.Chromosome.Author,
		DNA: new_dna})

	// @TODO hide interface
	p.Chain.Miner.Prove(new_bacteria.Chromosome)

	p.Bacts = append(p.Bacts, new_bacteria)

	return new_bacteria
}

func (p *Population) GetTrait(b *Bacteria, attitude_id string) uint {
	if p.Traits == nil {
		return 0
	}

	return b.Chromosome.DNA.MatchPatternCount(
		p.Traits[attitude_id].Pattern)
}

func (p *Population) Kill(target *Bacteria) {
	alive := make([]*Bacteria, 0)

	for _, b := range p.Bacts {
		if b != target {
			alive = append(alive, b)
		}
	}

	p.Bacts = alive
}

func (p *Population) DeliverChild() {
	newborn := p.Chain.Miner.GetMined()
	if newborn == nil {
		return
	}
	for _, b := range p.Bacts {
		if b.Chromosome == newborn {
			b.Born = true
			return
		}
	}
}

func (p *Population) String() string {
	return fmt.Sprintf("TRAITS:\n%s\nBACTS:\n%s\n", p.Traits, p.Bacts)
}

func (p *Population) Clone() *Population {
	newPop := &Population{
		Chain: p.Chain,
		Bacts: make([]*Bacteria, 0),
	}

	for _, b := range p.Bacts {
		newPop.Bacts = append(newPop.Bacts, b.Clone())
	}

	return newPop
}
