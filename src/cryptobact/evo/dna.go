package evo

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
)

var _ = log.Println

type DNA struct {
	Length  int
	Genes []*Gene
	SlicerSeed int
}

const GENE_MAX_LENGTH = 8

var GeneSlicer = rand.New(rand.NewSource(0))

func NewEmptyDNA() *DNA {
	return &DNA{
		Length: 0,
		Genes: make([]*Gene, 0),
		SlicerSeed: 0,
	}
}

func NewRandDNA(length int, slicerSeed int) *DNA {
	dna := NewEmptyDNA()
	dna.Length = length
	dna.SlicerSeed = slicerSeed

	for i := 0; i <= length; i += 1 {
		newGene := dna.Expand()
		newGene.Value = uint(rand.Intn((1 << uint(newGene.Length)) - 1))
	}

	return dna
}

func Crossover(a *DNA, b *DNA) *DNA {
	newDna := NewEmptyDNA()

	if b.Length > a.Length {
		newDna.Genes = b.Genes
		newDna.Length = b.Length
	} else {
		newDna.Genes = a.Genes
		newDna.Length = a.Length
	}

	for i, gene := range newDna.Genes {
		// dominant gene determined above
		choice := math.Min
		if gene.Length%2 == 0 {
			choice = math.Max
		}

		newDna.Genes[i] = &Gene{
			Value: uint(choice(float64(a.Genes[i].Value), float64(
				b.Genes[i].Value))),
			Length: gene.Length,
		}
	}

	return newDna
}

func (dna *DNA) Mutate(probability float64, rate float64) {
	for _, gene := range dna.Genes {
		if rand.Float64() < probability {
			gene.Value ^= uint(math.Min(math.Abs(rand.NormFloat64() * rate),
				float64((uint(1) << gene.Length) - 1)))
		}
	}
}

func (dna *DNA) Recombine(pattern string) {
	//@TODO
}

func (dna *DNA) GetNormGene(index int) float64 {
	return float64(dna.Genes[index].Value) / ((1 << (GENE_MAX_LENGTH - 1)) - 1)
}

func (dna *DNA) MatchPatternCount(pattern string) int {
	return 0 //@TODO
}

func (dna *DNA) Expand() *Gene {
	GeneSlicer.Seed(int64(dna.SlicerSeed))

	dna.SlicerSeed = GeneSlicer.Int()

	newGene := &Gene{
		Value: 0,
		Length: uint((dna.SlicerSeed % GENE_MAX_LENGTH) + 1),
	}

	dna.Genes = append(dna.Genes, newGene)

	return newGene
}

func (d *DNA) NumDiff(d2 *DNA) int {
	return 0 //@TODO
}

func (d *DNA) BitDiff(d2 *DNA) int {
	return 0 //@TODO
}

func (dna *DNA) String() string {
	geneStrs := make([]string, 0)
	for _, gene := range dna.Genes {
		geneStrs = append(geneStrs,
			fmt.Sprintf(fmt.Sprintf("%%0%db", gene.Length), gene.Value))
	}
	return strings.Join(geneStrs, "|")
}
