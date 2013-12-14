package infektor

import (
	"cryptobact/evo"
	"cryptobact/infektor/transport"

	"log"
)

type Infektor struct {
	amount    int
	rate      float64
	transport transport.Transporter
	ch        transport.InfectoChan
}

func NewInfektor(amount int, rate float64, t transport.Transporter) *Infektor {
	return &Infektor{transport: t, amount: amount, rate: rate}
}

func (ifk *Infektor) Serve() {
	ifk.ch = ifk.transport.Catch()
}

func (ifk *Infektor) Catch() *evo.Population {
	select {
	case pop2 := <-ifk.ch:
		return pop2
	default:
		return nil
	}
}

func (ifk *Infektor) Spread(pop *evo.Population) {
	bacts := make([]*evo.Bacteria, 0)

	for _, b := range pop.Bacts {
		if b.Born {
			bacts = append(bacts, b)
		}

		if len(bacts) >= int(float64(ifk.amount)*ifk.rate) {
			break
		}
	}

	if len(bacts) >= int(float64(ifk.amount)*ifk.rate) {
		log.Println("spreading infektion with amount of", ifk.amount)
		piligrims := pop.Splice(bacts[:ifk.amount])
		ifk.transport.Infect(piligrims)
	}
}
