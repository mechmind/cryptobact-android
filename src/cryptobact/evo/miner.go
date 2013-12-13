package evo

import (
	"fmt"
	"log"
	"math/big"
)

var _ = fmt.Print

type ControlCommand struct {
	command string
	task    *Chromosome
}

type ControlChan chan *ControlCommand

type TaskChan chan *Chromosome

type Miner struct {
	Difficulty *big.Int
	CtrlChan   ControlChan
	InChan     TaskChan
	OutChan    TaskChan
}

func NewMiner(difficulty uint, bufSize int) *Miner {
	threshold := big.NewInt(1)
	threshold.Lsh(threshold, difficulty)
	return &Miner{Difficulty: threshold,
		CtrlChan: make(ControlChan, bufSize),
		InChan:   make(TaskChan, bufSize),
		OutChan:  make(TaskChan, bufSize),
	}
}

func (m *Miner) Prove(chromo *Chromosome) {
	m.InChan <- chromo
}

func (m *Miner) GetMined() *Chromosome {
	select {
	case mined := <-m.OutChan:
		return mined
	default:
		return nil
	}
}

func (m *Miner) Start() {
	killChan := make(chan bool, 1)

	go mineManager(m.CtrlChan, m.InChan, killChan)
	go mineFacility(m.Difficulty, m.InChan, m.OutChan, killChan)
}

func mineManager(ctrlCh ControlChan, inCh TaskChan, killCh chan bool) {
	var work *Chromosome

	jobs := make([]*Chromosome, 0)

	for {
		cmd := <-ctrlCh
		switch cmd.command {
		case "new":
			jobs = append(jobs, cmd.task)
		case "cancel":
			if work.Hash(0) == cmd.task.Hash(0) {
				killCh <- true
			}

			for i, v := range jobs {
				if v.Hash(0) == cmd.task.Hash(0) {
					jobs = append(jobs[:i], jobs[i+1:]...)
					break
				}
			}
		case "getwork":
			work, jobs = jobs[0], jobs[:len(jobs)-1]
			inCh <- work
		default:
			log.Println("miner received incorrect command", cmd.command)
		}
	}
}

func mineFacility(diff *big.Int,
	inCh TaskChan, outCh TaskChan, killCh chan bool) {
	for {
		task := <-inCh

		nonce := uint(0)
		for {
			select {
			case <-killCh:
				log.Println("miner stopped at nonce", nonce, task)
				break
			default:
			}

			hash := task.Hash(nonce)
			if hash.Cmp(diff) <= 0 {
				task.Nonce = nonce
				task.CurrHash = hash
				outCh <- task
				break
			} else {
				nonce += 1
			}
		}
	}
}
