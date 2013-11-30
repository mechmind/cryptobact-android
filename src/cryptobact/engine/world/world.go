package world

import(
	"cryptobact/engine/bact"
	"cryptobact/engine/food"
	"cryptobact/engine/acid"
	"cryptobact/engine/clot"
)

const(
	FOOD_TICKS = 50
)

type World struct {
	Bacts []Bact
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
