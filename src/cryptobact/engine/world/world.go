package world

import "math/rand"
import . "cryptobact/evo/population"
import . "cryptobact/engine/food"
import . "cryptobact/engine/acid"
import . "cryptobact/engine/clot"

const (
	FOOD_TICKS = 50
	FOOD_PER_TICK = 3
)

type World struct {
	MyPopulation *Population
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
			w.Food = append(w.Food[:k], w.Food[k + 1:]...)
		}
	}
}
