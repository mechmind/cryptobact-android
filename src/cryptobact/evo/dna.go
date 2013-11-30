package evo

import "math"
import "math/big"
import "math/rand"
//import "regexp"
import "strings"

import "fmt"

type DNA struct {
    length int
    gene_pos []uint
    seq *big.Int
}

const GENE_SLICER_SEED = 0xDEADBEEE
const GENE_MAX_LENGTH = 8

var GeneSlicer = rand.New(rand.NewSource(0))

func NewEmptyDNA() *DNA {
    return &DNA{0, make([]uint, 0), big.NewInt(0)}
}

func NewRandDNA(length int) *DNA {
    dna := NewEmptyDNA()
    dna.length = length
    chunk := big.NewInt(0)
    for i := 0; i <= length; i += 64 {
        offset := uint(math.Min(64, float64(length - i)))
        chunk.SetInt64(rand.Int63() & (1 << offset - 1))
        dna.seq = dna.seq.Or(dna.seq, chunk.Lsh(chunk, uint(i)))
    }

    GeneSlicer.Seed(GENE_SLICER_SEED)
    for total := 0; total <= length; {
        gene_len := GeneSlicer.Intn(GENE_MAX_LENGTH) + 1
        dna.gene_pos = append(dna.gene_pos, uint(gene_len))
        total += gene_len
    }

    return dna
}

func Crossover(a *DNA, b *DNA) *DNA {
    a_genes := a.Genes()
    b_genes := b.Genes()

    new_dna := NewEmptyDNA()

    if len(b_genes) > len(a_genes) {
        a_genes, b_genes = b_genes, a_genes
        new_dna.gene_pos = b.gene_pos
        new_dna.length = b.length
    } else {
        new_dna.gene_pos = a.gene_pos
        new_dna.length = a.length
    }

    new_genome := big.NewInt(0)
    new_gene := big.NewInt(0)
    offset := uint(0)
    for i, gene_len := range new_dna.gene_pos {
        dominant := int64(0)
        if i % 2 == 0 {
            dominant = int64(math.Max(float64(a_genes[i]),
                float64(b_genes[i])))
        } else {
            dominant = int64(math.Min(float64(a_genes[i]),
                float64(b_genes[i])))
        }
        new_gene.SetInt64(int64(dominant))
        new_gene.Lsh(new_gene, uint(offset))
        new_genome.Or(new_genome, new_gene)
        offset += gene_len
    }

    new_dna.seq = new_genome

    return new_dna
}

func (dna *DNA) Mutate(probability float64, rate float64) {
    one := big.NewInt(1)
    bit_mask := big.NewInt(0)
    offset := uint(0)
    for _, gene_len := range dna.gene_pos {
        if rand.Float64() < probability {
            bits_to_change := uint(math.Abs(math.Floor(rand.NormFloat64() * rate)))
            bit_mask.Set(one)
            bit_mask.Lsh(bit_mask, uint(math.Min(
                float64(bits_to_change), float64(gene_len))))
            bit_mask.Sub(bit_mask, one)
            mutation := big.NewInt(rand.Int63())
            mutation.And(mutation, bit_mask)
            mutation.Lsh(mutation, offset)
            dna.seq.Xor(dna.seq, mutation)
        }
        offset += gene_len
    }
}

func (dna *DNA) Recombine(pattern string) {
    for i := 0; i < dna.length; i++ {
        find := -1
        pos := 0
        for j, ch := range pattern {
            pos = j
            if ch == '.' {
                continue;
            }

            if ch == '1' && dna.seq.Bit(i + j) != 1 {
                find = 1
                break
            }

            if ch == '0' && dna.seq.Bit(i + j) != 0 {
                find = 0
                break
            }
        }

        if find >= 0 {
            for k := i + pos; k < dna.length; k++ {
                if dna.seq.Bit(k) == uint(find) {
                    swap := dna.seq.Bit(k - 1)
                    dna.seq.SetBit(dna.seq, k - 1, uint(find))
                    dna.seq.SetBit(dna.seq, k, swap)
                    return;
                }
            }
        } else {
            i += pos - 1
        }
    }
}

func (dna *DNA) Genes() []uint {
    genes := make([]uint, 0)
    genome := big.NewInt(0)
    gene := big.NewInt(0)
    genome.Set(dna.seq)
    bit_mask := big.NewInt(0)
    one := big.NewInt(1)
    for i  := range dna.gene_pos {
        gene.Set(genome)
        bit_mask.SetInt64(1)
        bit_mask.Sub(bit_mask.Lsh(bit_mask, uint(dna.gene_pos[i])), one)
        gene.And(gene, bit_mask)
        genes = append(genes, uint(gene.Uint64()))
        genome.Rsh(genome, uint(dna.gene_pos[i]))
    }
    return genes
}

func (dna *DNA) MatchPatternCount(pattern string) uint {
    count := uint(0)
    for i := 0; i < dna.length; i++ {
        found := true
        for j, ch := range pattern {
            if ch == '.' {
                continue;
            }

            if ch == '1' && dna.seq.Bit(i + j) != 1 {
                found = false
                break
            }

            if ch == '0' && dna.seq.Bit(i + j) != 0 {
                found = false
                break
            }
        }

        if found {
            count += 1
        }
    }
    return count
}

func (dna *DNA) String() string {
    gene_strs := make([]string, 0)
    genes := dna.Genes()
    for i := len(genes) - 1; i >= 0; i-- {
        gene := genes[i]
        gene_strs = append(gene_strs,
            fmt.Sprintf(fmt.Sprintf("%%0%db", dna.gene_pos[i]), gene))
    }
    return strings.Join(gene_strs, "|")
}
