package engine

import . "cryptobact/evo/bacteria"
import . "cryptobact/engine/food"
import . "cryptobact/engine/acid"
import . "cryptobact/engine/clot"
import . "cryptobact/engine/world"
import . "cryptobact/engine/grid"
import . "cryptobact/engine/action"

const(
	WIDTH = 16
	HEIGHT = 24
)

type Updater interface {
	Update()
}

func Loop(updater Updater) {
	world := world.World{}
	chain := Chain{}

	grid := make(Grid, WIDTH)
	for x := 0; x < WIDTH; x++ {
		grid[x] = make([]Cell, HEIGHT)
	}

	world.Width = WIDTH
	world.Height = HEIGHT
	world.MyPopulation = NewPopulation(chain)

	// FIXME infinite loop goes here
	tick := 0
	if world.ShouldSpawnFood(tick) {
		world.SpawnFood()
		grid.CalcWeights(&world)
	}

	for _, bact := range world.MyPopulation.GetBacts() {
		a := action.GetAction(bact, &grid, &world)
		a.Apply(&world)
	}

	// FIXME call world.CleanFood()
	// FIXME call updater.Update(&world)

	return
}
