package engine

import (
	"cryptobact/evo"
	"cryptobact/infektor"
	"cryptobact/infektor/transport"

	"log"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
)

var _ = log.Println

const (
	WIDTH  = 16
	HEIGHT = 24

	FOOD_TICKS    = 500
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

	MAX_NANO = 999999999
)

type Updater interface {
	Update(*World)
}

func Loop(updater Updater) {
	runtime.GOMAXPROCS(2)

	defer func() {
		if err := recover(); err != nil {
			log.Println("work failed:", err)
			log.Printf("DEBUG: %s", debug.Stack())
		}
	}()

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
	targetTPS := TARGET_TPS
	for {
		startTime := time.Now()

		world.SpawnFood()

		for _, population := range world.Populations {
			if world.Notch(500) {
				log.Println(population)
			}

			SimulatePopulation(world, population)
		}

		world.CleanFood()

		updater.Update(world)

		if world.Notch(100) {
			infection := infektor.Catch()
			if infection != nil {
				exists := false
				for _, p := range world.Populations {
					if infection.Chain.Author == p.Chain.Author {
						exists = true
						break
					}
				}

				if !exists {
					log.Println("received infection size", len(infection.Bacts))
					log.Println(infection)
					infection.Chain.Miner = miner
					infection.Env = world.Populations[0].Env
					world.Populations = append(world.Populations, infection)
				}
			}
		}

		if world.Notch(110) {
			infektor.Spread(world.Populations[0])
		}

		world.Step()

		maxTPS := EstimateTPS(world.GetSmallTick(), startTime, &estTPSAvg)

		CalibrateTPS(maxTPS, targetTPS)
		CalibrateMiner(miner)

		realTPS := EstimateTPS(world.GetSmallTick(), startTime, &realTPSAvg)

		targetTPS = CorrectTargetTPS(realTPS, targetTPS)

		if world.Notch(500) {
			log.Printf("current miner hash rate is %.3f kh/s\n",
				miner.GetHashRate())
			log.Printf("current ticks per second is %d (%d) of %d max\n",
				realTPS, targetTPS, maxTPS)
			for i, p := range world.Populations {
				log.Printf("current population {%d} size is %d\n",
					i,
					len(p.Bacts))
			}
		}
	}
}

func EstimateTPS(smallTick int, startTime time.Time, TPSAvg *[]int) int {
	nano := time.Since(startTime).Nanoseconds()
	if nano == 0 {
		return -1
	}

	TPS := int(MAX_NANO / nano)
	(*TPSAvg)[smallTick%TPS_WINDOW_SIZE] = TPS

	sum := 0
	for _, v := range *TPSAvg {
		sum += v
	}

	return sum / TPS_WINDOW_SIZE
}

func SimulatePopulation(world *World, population *evo.Population) {
	for _, bact := range population.Bacts {
		a := GetAction(population, bact, world)
		if a != nil {
			a.Apply()
		} else {
			log.Println("leave me here comrades", bact)
		}
	}

	population.DeliverChild()
	world.GetOld(population)
	world.ApplyAcid(population)
}

func InitPopulation(world *World, population *evo.Population) {
	for _, b := range population.Bacts {
		b.X = rand.Float64() * float64(world.Width)
		b.Y = rand.Float64() * float64(world.Height)
	}
}

func CalibrateMiner(m *evo.Miner) {
	if m.GetHashRate() <= 0 {
		return
	}

	newDiff := int(MINER_BASE_DIFF - math.Ceil(
		float64(m.GetHashRate()/MINER_BASE_RATE)) + 1)

	if m.Difficulty != newDiff {
		log.Printf("miner calibration %d -> %d", m.Difficulty, newDiff)
	}

	m.SetDifficulty(newDiff)
}

func CalibrateTPS(rate int, target int) {
	if rate <= target {
		return
	}

	avg := 1.0 / float64(rate) * MAX_NANO

	sleep := math.Floor(float64(rate)/float64(target)) * avg

	time.Sleep(time.Duration(sleep) * time.Nanosecond)
}

func CorrectTargetTPS(realTPS int, target int) int {
	if math.Abs(float64(TARGET_TPS-realTPS)) <= 10 {
		return target + int(math.Ceil(float64(TARGET_TPS-realTPS)/2))
	} else {
		return target
	}
}
