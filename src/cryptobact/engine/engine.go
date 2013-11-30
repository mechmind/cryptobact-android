package engine

import "runtime"
//import "fmt"

import "cryptobact/evo"

const(
    WIDTH = 16
    HEIGHT = 24
)

type Updater interface {
    Update()
}

func Loop(updater Updater) {
    runtime.GOMAXPROCS(2)

    miner := evo.NewMiner(146)
    miner.Start()

    chain := &evo.Chromochain{}

    world := &World{}

    options := &evo.PopulationOptions{
        Attitudes: map[string]*evo.Attitude{
            "lust": &evo.Attitude{"111.1", 4},
            "glut": &evo.Attitude{"10101", 2},
        },
        MutateProbability: 0.5,
        MutateRate: 1,
        RecombinationChance: 1.0,
        RecombinationDrop: 10,
    }

    grid := make(Grid, WIDTH)
    for x := 0; x < WIDTH; x++ {
        grid[x] = make([]Cell, HEIGHT)
    }

    world.Width = WIDTH
    world.Height = HEIGHT
    world.MyPopulation = evo.NewPopulation(miner, chain, options)

	for _, b := range world.MyPopulation.GetBacts() {
		b.TTL = int(10000 * float64(world.MyPopulation.GetGene(b, 7)) / 10)
		b.Energy = 1000 * float64(world.MyPopulation.GetGene(b, 11)) / 10
	}

    tick := 0
    for {
        world.SpawnFood(tick)
        grid.CalcWeights(world)

        for _, bact := range world.MyPopulation.GetBacts() {
            a := GetAction(bact, &grid, world)
            a.Apply(bact, world)
        }

        world.MyPopulation.CatchNewBorn()
        //fmt.Println(world.MyPopulation.GetBacts())
        world.CleanFood()
		world.GetOld()
        // FIXME call updater.Update(&world)

        tick += 1
    }
}
