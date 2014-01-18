package evo

import (
	"fmt"
	"log"
)

var _ = fmt.Print
var _ = log.Print

type Creature interface {
	Reproduce(a *Creature) *Creature
	GetChromosome() *Chromosome
}

type Population struct {
	// @TODO author
	Traits TraitMap
	Chain  *Chromochain

	Creatures []*Creature
}
