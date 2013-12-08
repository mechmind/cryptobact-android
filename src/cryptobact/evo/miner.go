package evo

import "math/big"
import "fmt"

var _ = fmt.Print

type Miner struct {
	Difficulty *big.Int
	InChan     chan *Chromosome
	OutChan    chan *Chromosome
}

func NewMiner(difficulty uint) *Miner {
	threshold := big.NewInt(1)
	threshold.Lsh(threshold, difficulty)
	return &Miner{Difficulty: threshold,
		InChan:  make(chan *Chromosome, 256),
		OutChan: make(chan *Chromosome, 256),
	}
}

func (m *Miner) Mine(chromo *Chromosome) {
	m.InChan <- chromo
}

func (m *Miner) Extract() *Chromosome {
	select {
	case newborn := <-m.OutChan:
		return newborn
	default:
		return nil
	}
}

func (m *Miner) Start() {
	go func(m *Miner) {
		for {
			task := <-m.InChan

			nonce := uint(0)
			for {
				hash := task.Hash(nonce)
				if hash.Cmp(m.Difficulty) <= 0 {
					task.Nonce = nonce
					task.CurrHash = hash
					m.OutChan <- task
					break
				} else {
					nonce += 1
				}
			}
		}
	}(m)
}
