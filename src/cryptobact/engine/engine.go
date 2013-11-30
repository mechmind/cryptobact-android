package engine

import . "cryptobact/evo/bacteria"
import . "cryptobact/engine/food"
import . "cryptobact/engine/acid"
import . "cryptobact/engine/clot"
import . "cryptobact/engine/world"
import . "cryptobact/engine/grid"

func Loop() {
	world := world.World{}
	grid := grid.Grid{}
    chain := Chain{}

	world.MyPopulation = NewPopulation(chain)

	// FIXME infinite loop goes here
	tick := 0
	if world.ShouldSpawnFood(tick) {
		world.SpawnFood()
		grid.CalcWeights(&world)
	}

	for _, bact := range world.MyPopulation.Bacts {
		action := bact.GetAction(grid, world)
		action.Apply(&world)
	}

	// FIXME call render.Update()

	return
}
