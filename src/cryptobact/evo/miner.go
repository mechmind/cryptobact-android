package evo

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"math/rand"
)

var _ = fmt.Print

type Task struct {
	reldiff int
	chromo  *Chromosome
}

type Miner struct {
	Difficulty int
	Threshold  *big.Int
	nonce      int
	khs        float64
	getwork    chan bool
	kill       chan bool
	cancel     chan *Chromosome
	task       chan Task
	sendwork   chan Task
	proved     chan *Chromosome
	started    chan *Chromosome
}

func NewMiner(difficulty int, bufSize int) *Miner {
	m := &Miner{
		getwork:  make(chan bool, bufSize),
		kill:     make(chan bool, bufSize),
		cancel:   make(chan *Chromosome, bufSize),
		task:     make(chan Task, bufSize),
		sendwork: make(chan Task, bufSize),
		proved:   make(chan *Chromosome, bufSize),
		started:  make(chan *Chromosome, bufSize),
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

func (m *Miner) Prove(chromo *Chromosome, reldiff int) {
	m.task <- Task{reldiff, chromo}
}

func (m *Miner) Cancel(chromo *Chromosome) {
	m.cancel <- chromo
}

func (m *Miner) GetStarted() *Chromosome {
	select {
	case started := <-m.started:
		return started
	default:
		return nil
	}
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
	var work Task

	jobs := make([]Task, 0)

	for {
		select {
		case t := <-m.task:
			log.Printf("miner new task %s, reldiff %d", t.chromo, t.reldiff)
			jobs = append(jobs, t)
		case <-m.getwork:
			if len(jobs) == 0 {
				log.Println("miner getwork and no work avail")
				work = <-m.task
				log.Println("miner received first task", work)
				m.sendwork <- work
			} else {
				log.Println("miner getwork")
				pos := rand.Intn(len(jobs))
				work, jobs = jobs[pos], append(jobs[:pos], jobs[pos+1:]...)
				m.sendwork <- work
			}
		case t := <-m.cancel:
			log.Println("miner cancel", t)
			for i, v := range jobs {
				if v.chromo == t {
					jobs = append(jobs[:i], jobs[i+1:]...)
					break
				}
			}

			if work.chromo == t {
				m.kill <- true
			}
		}
	}
}

func mineFacility(m *Miner) {
	for {
		m.getwork <- true
		task := <-m.sendwork
		//m.started <- task.chromo
		chromo := task.chromo

		difficulty := m.Difficulty

		if task.reldiff <= 2 {
			difficulty -= 1
		}

		if task.reldiff > 2 && task.reldiff <= 5 {
			difficulty += 3
		}

		threshold := big.NewInt(1)
		threshold.Lsh(threshold, uint(difficulty))

		log.Printf("miner start mining at diff %020x\n", threshold)
		startTime := time.Now()
		measureTime := time.Now()
		m.nonce = 0
		m.khs = 0
		nonce := 0
	outerFor:
		for {
			select {
			case <-m.kill:
				log.Println("miner killed at nonce", nonce, chromo)
				break outerFor
			default:
			}

			if time.Since(measureTime) > 2*time.Second {
				m.khs = float64(nonce-m.nonce) /
					float64(time.Since(measureTime).Seconds()) /
					1000.0

				measureTime = time.Now()

				m.nonce = nonce
			}

			hash := chromo.Hash(nonce)
			if hash.Cmp(threshold) <= 0 {
				log.Printf("miner successfully mined task at nonce %d, time %.2f sec",
					nonce,
					time.Since(startTime).Seconds())
				chromo.Nonce = nonce
				chromo.CurrHash = hash
				m.proved <- chromo
				break outerFor
			} else {
				nonce += 1
			}
		}
	}
}
