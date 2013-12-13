package engine

import (
	"cryptobact/evo"

	"math/big"
	"math/rand"
)

type World struct {
	Populations []*evo.Population
	Food        []*Food
	Acid        []*Acid
	Clot        []*Clot
	Width       int
	Height      int
	FoodTicks   int
	FoodPerTick int
	Tick        *big.Int
}

func NewWorld() *World {
	return &World{Tick: big.NewInt(0)}
}

func (w *World) SpawnFood() {
	if !w.Notch(w.FoodTicks) {
		return
	}

	for i := 0; i < w.FoodPerTick; i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		w.Food = append(w.Food, &Food{x, y, false})
	}

	return
}

func (w *World) CleanFood() {
	for k, f := range w.Food {
		if f.Eaten {
			if k+1 >= len(w.Food) {
				w.Food = w.Food[:k]
			} else {
				w.Food = append(w.Food[:k], w.Food[k+1:]...)
			}
		}
	}
}

func (w *World) GetOld(population *evo.Population) {
	for _, b := range population.Bacts {
		if !b.Born {
			continue
		}
		b.TTL -= int(population.GetGene(b, 17)/10.0 + 1)
	}
}

func (w *World) Step() {
	w.Tick = w.Tick.Add(w.Tick, big.NewInt(1))
}

func (w *World) Notch(notch int) bool {
	if w.Tick.Mod(w.Tick, big.NewInt(int64(notch))) == big.NewInt(0) {
		return true
	} else {
		return false
	}
}
