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
	WIDTH  = 16
	HEIGHT = 24

	FOOD_TICKS    = 20
	FOOD_PER_TICK = 10

	MINER_BASE_DIFF = 145
	MINER_BASE_RATE = 6.0
	MINER_BUFFER    = 255

	INFECT_WITH_SIZE  = 4
	INFECT_MULTIPLIER = 2

	INFECT_PORT_1 = 1234
	INFECT_PORT_2 = 2345
	INFECT_PORT_3 = 3456

	TARGET_TPS = 100

	CALIBRATE_MS = 200

	TPS_WINDOW_SIZE = 10
)

type Updater interface {
	Update(*World)
}

func Loop(updater Updater) {
	runtime.GOMAXPROCS(2)

	rand.Seed(time.Now().UnixNano())

	miner := evo.NewMiner(MINER_BASE_DIFF, MINER_BUFFER)
	miner.Start()

	chain := &evo.Chromochain{
		Author: uint64(rand.Int63()),
		Miner:  miner}

	world := NewWorld()

	traits := evo.TraitMap{
		"lust": &evo.Trait{"111.1", 4},
		"glut": &evo.Trait{"10101", 2},
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

	world.SpawnAcid()
	world.SpawnClot()
	//startTick := 0
	estTPSAvg := make([]int, TPS_WINDOW_SIZE)
	realTPSAvg := make([]int, TPS_WINDOW_SIZE)
	for {
		startTime := time.Now()

		world.SpawnFood()

		for _, population := range world.Populations {
			SimulatePopulation(world, population)
		}

		world.CleanFood()

		updater.Update(world)

		if world.Notch(100) {
			infection := infektor.Catch()
			if infection != nil {
				if infection.Chain.Author != chain.Author {
					log.Println("received infection size", len(infection.Bacts))
					log.Println(infection)
				}
			}
		}

		if world.Notch(110) {
			infektor.Spread(world.Populations[0])
		}

		world.Step()

		maxTPS := EstimateTPS(world.GetSmallTick(), startTime, &estTPSAvg)

		CalibrateTPS(maxTPS)
		CalibrateMiner(miner)

		realTPS := EstimateTPS(world.GetSmallTick(), startTime, &realTPSAvg)

		if world.Notch(500) {
			log.Printf("current miner hash rate is %.3f kh/s\n",
				miner.GetHashRate())
			log.Printf("current ticks per second is %d of %d max\n",
				realTPS, maxTPS)
			log.Printf("current My Population size is %d\n",
				len(world.Populations[0].Bacts))
		}
	}
}

func EstimateTPS(smallTick int, startTime time.Time, TPSAvg *[]int) int {
	nano := time.Since(startTime).Nanoseconds()
	TPS := int(999999999 / nano)
	(*TPSAvg)[smallTick%TPS_WINDOW_SIZE] = TPS

	sum := 0
	for _, v := range *TPSAvg {
		sum += v
	}

	return sum / TPS_WINDOW_SIZE
}

func SimulatePopulation(world *World, population *evo.Population) {
	for _, bact := range population.Bacts {
		if !bact.Born {
			continue
		}
		a := GetAction(population, bact, world)
		a.Apply()
	}

	population.DeliverChild()
	world.GetOld(population)
	world.ApplyAcid(population)
	world.ApplyClot(population)
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

func CalibrateMiner(m *evo.Miner) {
	if m.GetHashRate() <= 0 {
		return
	}

	m.SetThreshold(MINER_BASE_DIFF - int(m.GetHashRate()/MINER_BASE_RATE))
}

func CalibrateTPS(rate int) {
	if rate <= TARGET_TPS {
		return
	}

	avg := 1.0 / float64(rate) * 999999999

	sleep := (float64(rate)/float64(TARGET_TPS) - 1) * avg

	time.Sleep(time.Duration(sleep) * time.Nanosecond)
}
