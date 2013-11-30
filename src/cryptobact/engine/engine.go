package engine

import(
	"cryptobact/engine/api"
	"cryptobact/engine/bact"
	"cryptobact/engine/food"
	"cryptobact/engine/acid"
	"cryptobact/engine/clot"
	"cryptobact/engine/world"
	"cryptobact/engine/grid"
)

func Loop() {
	world := world.World{}
	grid := grid.Grid{}
	bacts := api.GetBacts(1)

	world.bacts = bacts

	// FIXME infinite loop goes here
	tick := 0
	if world.ShouldSpawnFood(tick) {
		world.SpawnFood()
		grid.CalcWeights(&world)
	}

	for i := range world.bacts {
		bact = world.bacts[i]
		action := bact.GetAction(grid, world)
		action.Apply(&world)
	}

	// FIXME call render.Update()

	return
}
