package evo

import (
	"fmt"
	"log"
	"math/big"
	"time"
)

var _ = fmt.Print

type Miner struct {
	Difficulty int
	Threshold  *big.Int
	nonce      int
	khs        float64
	getwork    chan bool
	kill       chan bool
	cancel     chan *Chromosome
	task       chan *Chromosome
	sendwork   chan *Chromosome
	proved     chan *Chromosome
}

func NewMiner(difficulty int, bufSize int) *Miner {
	m := &Miner{
		getwork:  make(chan bool, bufSize),
		kill:     make(chan bool, bufSize),
		cancel:   make(chan *Chromosome, bufSize),
		task:     make(chan *Chromosome, bufSize),
		sendwork: make(chan *Chromosome, bufSize),
		proved:   make(chan *Chromosome, bufSize),
	}

	m.SetDifficulty(difficulty)

	return m
}

func (m *Miner) SetDifficulty(difficulty int) {
	m.Difficulty = difficulty
	threshold := big.NewInt(1)
	threshold.Lsh(threshold, uint(difficulty))
	m.Threshold = threshold
}

func (m *Miner) Prove(chromo *Chromosome) {
	m.task <- chromo
}

func (m *Miner) Cancel(chromo *Chromosome) {
	m.cancel <- chromo
}

func (m *Miner) GetMined() *Chromosome {
	select {
	case mined := <-m.proved:
		return mined
	default:
		return nil
	}
}

func (m *Miner) Start() {
	go mineManager(m)
	go mineFacility(m)
}

func (m *Miner) GetHashRate() float64 {
	return m.khs
}

func mineManager(m *Miner) {
	var work *Chromosome

	jobs := make([]*Chromosome, 0)

	for {
		select {
		case t := <-m.task:
			//log.Println("miner new task", t)
			jobs = append(jobs, t)
		case <-m.getwork:
			if len(jobs) == 0 {
				log.Println("miner getwork and no work avail")
				work = <-m.task
				log.Println("miner received first task", work)
				m.sendwork <- work
			} else {
				log.Println("miner getwork")
				work, jobs = jobs[0], jobs[1:]
				m.sendwork <- work
			}
		case t := <-m.cancel:
			log.Println("miner cancel", t)
			for i, v := range jobs {
				if v == t {
					jobs = append(jobs[:i], jobs[i+1:]...)
					break
				}
			}

			if work == t {
				m.kill <- true
			}
		}
	}
}

func mineFacility(m *Miner) {
	for {
		m.getwork <- true
		task := <-m.sendwork

		log.Printf("miner start mining at diff %020x\n", m.Threshold)
		startTime := time.Now()
		measureTime := time.Now()
		m.nonce = 0
		m.khs = 0
		nonce := 0
		for {
			select {
			case <-m.kill:
				log.Println("miner killed at nonce", nonce, task)
				break
			default:
			}

			if time.Since(measureTime) > 2*time.Second {
				m.khs = float64(nonce-m.nonce) /
					float64(time.Since(measureTime).Seconds()) /
					1000.0

				measureTime = time.Now()

				m.nonce = nonce
			}

			hash := task.Hash(nonce)
			if hash.Cmp(m.Threshold) <= 0 {
				log.Printf("miner successfully mined task at nonce %d, time %.2f sec",
					nonce,
					time.Since(startTime).Seconds())
				task.Nonce = nonce
				task.CurrHash = hash
				m.proved <- task
				break
			} else {
				nonce += 1
			}
		}
	}
}
