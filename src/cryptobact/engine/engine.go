package engine

import (
	"cryptobact/evo"
	"cryptobact/infektor"
	"cryptobact/infektor/transport"

	"log"
	"math/rand"
	"runtime"
	"time"
)

var _ = log.Println

const (
	WIDTH         = 16
	HEIGHT        = 24
	FOOD_TICKS    = 20
	FOOD_PER_TICK = 10

	MINER_DIFF   = 149
	MINER_BUFFER = 255

	INFECT_WITH_SIZE  = 4
	INFECT_MULTIPLIER = 2

	INFECT_PORT_1 = 1234
	INFECT_PORT_2 = 2345
	INFECT_PORT_3 = 3456
)

type Updater interface {
	Update(*World)
}

func Loop(updater Updater) {
	runtime.GOMAXPROCS(2)

	rand.Seed(time.Now().UnixNano())

	miner := evo.NewMiner(MINER_DIFF, MINER_BUFFER)
	miner.Start()

	chain := &evo.Chromochain{
		Author: uint64(rand.Int63()),
		Miner:  miner}

	world := &World{}

	traits := evo.TraitMap{
		"lust": &evo.Trait{"111.1", 4},
		"glut": &evo.Trait{"10101", 2},
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
		evo.NewPopulation(chain, traits, nil),
	}

	infektor := infektor.NewInfektor(INFECT_WITH_SIZE, INFECT_MULTIPLIER,
		transport.NewUDP([]int{
			INFECT_PORT_1,
			INFECT_PORT_2,
			INFECT_PORT_3,
		}))

	InitPopulation(world, world.Populations[0])

	infektor.Serve()

	world.Tick = 0
	for {
		world.SpawnFood(world.Tick)
		grid.CalcWeights(world)

		for _, population := range world.Populations {
			SimulatePopulation(&grid, world, population)
		}

		world.CleanFood()
		updater.Update(world)

		if world.Tick%100 == 0 {
			infection := infektor.Catch()
			if infection != nil {
				if infection.Chain.Author != chain.Author {
					log.Println("received infection size", len(infection.Bacts))
					log.Println(infection)
				}
			}
		}

		if world.Tick%110 == 0 {
			infektor.Spread(world.Populations[0])
		}

		if world.Tick == 999 {
			world.Tick = 0
		} else {
			world.Tick += 1
		}
	}
}

func SimulatePopulation(grid *Grid, world *World, population *evo.Population) {
	for _, bact := range population.Bacts {
		if !bact.Born {
			continue
		}
		a := GetAction(population, bact, grid, world)
		a.Apply()
	}

	population.DeliverChild()
	world.GetOld(population)
}

func InitPopulation(world *World, population *evo.Population) {
	for _, b := range population.Bacts {
		b.X = rand.Float64() * float64(world.Width)
		b.Y = rand.Float64() * float64(world.Height)
		b.TTL = int(10000 * float64(population.GetGene(b, 7)) / 10)
		b.Energy = 1000 * float64(population.GetGene(b, 11)) / 10
		b.RotationSpeed = 10.0 + float64(population.GetGene(b, 4)/20)
	}
}
