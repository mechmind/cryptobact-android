package engine

import "math/rand"
import "cryptobact/evo"

const (
	FOOD_TICKS = 30
	FOOD_PER_TICK = 10
)

type World struct {
	MyPopulation *evo.Population
	Food []*Food
	Acid []*Acid
	Clot []*Clot
	Width int
	Height int
}

func (w *World) SpawnFood(tick int) {
	if (tick % FOOD_TICKS) != 0 {
		return
	}

	for i := 0; i < FOOD_PER_TICK; i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		w.Food = append(w.Food, &Food{x, y, false})
	}

	return
}

func (w *World) CleanFood() {
	for k, f := range w.Food {
		if f.Eaten {
			if k + 1 >= len(w.Food) {
				w.Food = w.Food[:k]
			} else {
				w.Food = append(w.Food[:k], w.Food[k + 1:]...)
			}
		}
	}
}

func (w *World) GetOld() {
	for _, b := range w.MyPopulation.GetBacts() {
		b.TTL -= int(w.MyPopulation.GetGene(b, 17) / 10.0 + 1)
	}
}
