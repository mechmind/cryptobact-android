package world

import . "cryptobact/evo/population"
import . "cryptobact/engine/food"
import . "cryptobact/engine/acid"
import . "cryptobact/engine/clot"

const FOOD_TICKS = 50

type World struct {
	MyPopulation *Population
	Food []Food
	Acid []Acid
	Clot []Clot
}

func ShouldSpawnFood(tick int) bool {
}

func (w *World) SpawnFood() {
	if (tick % FOOD_TICKS) != 0 {
		return
	}

	// TODO implement

	return
}
