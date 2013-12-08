package engine

import "math/rand"
import "cryptobact/evo"

type World struct {
	Populations []*evo.Population
	Food        []*Food
	Acid        []*Acid
	Clot        []*Clot
	Width       int
	Height      int
	FoodTicks   int
	FoodPerTick int
	Tick        int
}

func (w *World) SpawnFood(tick int) {
	if (tick % w.FoodTicks) != 0 {
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
	for _, b := range population.GetBacts() {
		if !b.Born {
			continue
		}
		b.TTL -= int(population.GetGene(b, 17)/10.0 + 1)
	}
}
