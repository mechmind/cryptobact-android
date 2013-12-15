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

	FOOD_TICKS    = 100
	FOOD_PER_TICK = 10

	MINER_BASE_DIFF = 146
	MINER_BASE_RATE = 6.0
	MINER_BUFFER    = 255

	INFECT_WITH_SIZE  = 3
	INFECT_MULTIPLIER = 3

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
			panic(err)
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
		"lust": &evo.Trait{"111.111", 4},
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

	//eTotal := int64(0)
	//eSpawnFood := int64(0)
	//ePopulationTotal := int0
	//ePopulationSpread := 0
	//eCounter := 0
	//

	var (
		eCounter          int64 = 0
		eTotal            int64 = 0
		eSpawnFood        int64 = 0
		ePopulationTotal  int64 = 0
		ePopulationSpread int64 = 0
		eCleanFood        int64 = 0
		eUpdater          int64 = 0
		eSimulate         int64 = 0
		eNewbornProcess   int64 = 0
		eCatch            int64 = 0
		eCalibration      int64 = 0
	)

	for {
		eCounter += 1
		tTotal := time.Now()
		//startTime := time.Now()

		tSpawnFood := time.Now()
		world.SpawnFood()
		eSpawnFood += time.Since(tSpawnFood).Nanoseconds()

		newborn := miner.GetMined()

		tPopulationTotal := time.Now()
		for _, population := range world.Populations {
			if world.Notch(1000) {
				log.Println(population)
			}

			if world.Notch(110) {
				tPopulationSpread := time.Now()
				infektor.Spread(population)
				ePopulationSpread += time.Since(tPopulationSpread).Nanoseconds()
			}

			tSimulate := time.Now()
			SimulatePopulation(world, population)
			eSimulate += time.Since(tSimulate).Nanoseconds()

			tNewbornProcess := time.Now()
			if newborn != nil && newborn.Author == population.Chain.Author {
				for _, b := range population.Bacts {
					if b.Chromosome == newborn {
						b.RenewTTL()
						b.Born = true
					}
				}
			}
			eNewbornProcess += time.Since(tNewbornProcess).Nanoseconds()
		}
		ePopulationTotal += time.Since(tPopulationTotal).Nanoseconds()

		tCleanFood := time.Now()
		world.CleanFood()
		eCleanFood += time.Since(tCleanFood).Nanoseconds()

		tUpdater := time.Now()
		updater.Update(world)
		eUpdater += time.Since(tUpdater).Nanoseconds()

		if world.Notch(100) {
			tCatch := time.Now()
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
			eCatch += time.Since(tCatch).Nanoseconds()
		}

		world.Step()

		eTotal += time.Since(tTotal).Nanoseconds()

		tCalibration := time.Now()
		maxTPS := EstimateTPS(world.GetSmallTick(), tTotal, &estTPSAvg)

		CalibrateTPS(maxTPS, targetTPS)
		CalibrateMiner(miner)

		realTPS := EstimateTPS(world.GetSmallTick(), tTotal, &realTPSAvg)
		eCalibration += time.Since(tCalibration).Nanoseconds()

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

			log.Printf("aggregated profile for %d ticks", eCounter)

			lCounter := float64(eCounter) * 1.0e9
			log.Printf("total time       : %.4f", float64(eTotal)/lCounter)
			log.Printf("spawn food       : %.4f", float64(eSpawnFood)/lCounter)
			log.Printf("population total : %.4f", float64(ePopulationTotal)/lCounter)
			log.Printf("update world     : %.4f", float64(eUpdater)/lCounter)
			log.Printf("clean food       : %.4f", float64(eCleanFood)/lCounter)
			log.Printf("population spread: %.4f", float64(ePopulationSpread)/lCounter)
			log.Printf("simulate         : %.4f", float64(eSimulate)/lCounter)
			log.Printf("newborn process  : %.4f", float64(eNewbornProcess)/lCounter)
			log.Printf("catch            : %.4f", float64(eCatch)/lCounter)
			log.Printf("calibration      : %.4f", float64(eCalibration)/lCounter)

			eCounter = 0
			eTotal = 0
			eSpawnFood = 0
			ePopulationTotal = 0
			eUpdater = 0
			eCleanFood = 0
			ePopulationSpread = 0
			eSimulate = 0
			eNewbornProcess = 0
			eCatch = 0
			eCalibration = 0
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
		a := getCachedAction(world, bact)
		if a == nil {
			a = GetAction(population, bact, world)
			if a == nil {
				continue
			}
			switch a.(type) {
			case *ActionMove:
				cacheAction(world, bact, a)
			}
		}
		a.Apply()
	}

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

func cacheAction(w *World, b *evo.Bacteria, a Action) {
	// FIXME replace hardcoded ttl
	ttl := 50
	for _, c := range w.ActionCache {
		if c.B == b {
			c.Action = a
			c.TTL = ttl
			return
		}
	}
	action := &ActionCache{b, a, ttl}
	w.ActionCache = append(w.ActionCache, action)
}

func getCachedAction(w *World, b *evo.Bacteria) Action {
	return nil
	inertia := false
	if dist(0, 0, b.Inertia.X, b.Inertia.Y) > 1e-5 {
		inertia = true
	}
	for k, c := range w.ActionCache {
		if c.B == b {
			if c.TTL <= 1 || inertia {
				if k+1 >= len(w.ActionCache) {
					w.ActionCache = w.ActionCache[:k]
				} else {
					w.ActionCache = append(w.ActionCache[:k], w.ActionCache[k+1:]...)
				}
			}
			if inertia {
				return nil
			}
			c.TTL -= 1
			return c.Action
		}
	}
	return nil
}
