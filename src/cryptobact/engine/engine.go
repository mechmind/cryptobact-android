package engine

import "runtime"
import "math/rand"
import "time"
import "log"

import "cryptobact/evo"
import "cryptobact/infektor"

const(
    WIDTH = 16
    HEIGHT = 24
	FOOD_TICKS = 20
	FOOD_PER_TICK = 10
)

var Miner *evo.Miner

type Updater interface {
    Update(*World)
}

func Loop(updater Updater) {
    runtime.GOMAXPROCS(2)

    rand.Seed(time.Now().UnixNano())

    Miner = evo.NewMiner(149)
    Miner.Start()

    chain := &evo.Chromochain{Author: uint64(rand.Int63())}

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
	world.FoodTicks = FOOD_TICKS
	world.FoodPerTick = FOOD_PER_TICK
    world.Populations = []*evo.Population{
        evo.NewPopulation(Miner, chain, options),
    }

    infektor := infektor.NewInfektor([]uint{
            2345,
            4567,
            5678,
        },
    )

    InitPopulation(world, world.Populations[0])

    infections := infektor.Listen()
    infektor.Spread(world.Populations[0], 1 * time.Second)

    world.Tick = 0
    for {
        world.SpawnFood(world.Tick)
        grid.CalcWeights(world)

        for _, population := range world.Populations {
            SimulatePopulation(&grid, world, population)
        }

        ProcessInfections(world, infections)

        world.CleanFood()
        updater.Update(world)

		if world.Tick == 999 {
			world.Tick = 0
		} else {
			world.Tick += 1
		}
    }
}

func SimulatePopulation(grid *Grid, world *World, population *evo.Population) {
    for _, bact := range population.GetBacts() {
        if !bact.Born {
            continue
        }
        a := GetAction(population, bact, grid, world)
        a.Apply()
    }

    population.CatchNewBorn()
    world.GetOld(population)
}

func InitPopulation(world *World, population *evo.Population) {
	for _, b := range population.GetBacts() {
		b.X = rand.Float64() * float64(world.Width)
		b.Y = rand.Float64() * float64(world.Height)
		b.TTL = int(10000 * float64(population.GetGene(b, 7)) / 10)
		b.Energy = 1000 * float64(population.GetGene(b, 11)) / 10
		b.RotationSpeed = 10.0 + float64(population.GetGene(b, 4) / 20)
	}
}

func ProcessInfections(world *World, infections chan *evo.Chromosome) {
    stop_for := false
    for !stop_for {
        select {
        case new_chromo := <- infections:
            new_chain := &evo.Chromochain{
                Author: new_chromo.Author,
                Initial: new_chromo}

            skip := false
            for _, p := range world.Populations {
                if p.Chain.Author == new_chromo.Author {
                    skip = true
                    break
                }
            }

            if (skip) {
                continue
            }

            log.Println("INFECT", new_chain.Author)

            new_population := evo.NewPopulation(Miner, new_chain, nil)
            InitPopulation(world, new_population)

            world.Populations = append(world.Populations, new_population)
        default:
            stop_for = true
        }
    }
}
