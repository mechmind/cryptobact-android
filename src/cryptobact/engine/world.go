package engine

import (
	"cryptobact/evo"

	"math"
	"math/big"
	"math/rand"
)

const (
	MAX_FOODS         = 100
	MAX_ACID_SPOTS    = 5
	MAX_ACID_CON      = 0.5
	MAX_CLOT_SPOTS    = 5
	MAX_CLOT_DENS     = 0.5
	MIN_FOOD_CALORIES = 30.0
	MAX_FOOD_CALORIES = 50.0
)

type World struct {
	Populations []*evo.Population
	Food        []*Food
	Acid        []*Acid
	Clot        []*Clot
	Width       int
	Height      int
	FoodTicks   int
	FoodPerTick int
	Tick        *big.Int
}

func NewWorld() *World {
	return &World{
		Populations: make([]*evo.Population, 0),
		Tick:        big.NewInt(0),
	}
}

// adds some acid spots
func (w *World) SpawnAcid() {
	for i := 0; i < rand.Intn(MAX_ACID_SPOTS); i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		concentration := rand.Float64() * MAX_ACID_CON
		w.Acid = append(w.Acid, &Acid{x, y, concentration})
	}
}

// adds some clots
func (w *World) SpawnClot() {
	for i := 0; i < rand.Intn(MAX_CLOT_SPOTS); i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		density := rand.Float64() * MAX_CLOT_DENS
		w.Clot = append(w.Clot, &Clot{x, y, density})
	}
}

// randomly spawns food
func (w *World) SpawnFood() {
	if !w.Notch(w.FoodTicks) {
		return
	}

	if len(w.Food) > MAX_FOODS {
		return
	}

	for i := 0; i < w.FoodPerTick; i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		min := MIN_FOOD_CALORIES
		max := MAX_FOOD_CALORIES
		r := rand.Float64()
		calories := min + r*(max-min)
		w.Food = append(w.Food, &Food{x, y, calories, false})
	}
}

// removes eaten food from the map
func (w *World) CleanFood() {
	for k, f := range w.Food {
		if f.Eaten {
			if k+1 >= len(w.Food) {
				w.Food = w.Food[:k]
			} else {
				w.Food = append(w.Food[:k], w.Food[k+1:]...)
			}
		}
	}
}

// makes bacteries a little older
func (w *World) GetOld(population *evo.Population) {
	for _, b := range population.Bacts {
		//if b.Born || b.Labouring || true {
		b.TTL -= 1
		//}
	}
}

func (w *World) Step() {
	w.Tick = w.Tick.Add(w.Tick, big.NewInt(1))
}

func (w *World) Notch(notch int) bool {
	if w.GetSmallTick()%notch == 0 {
		return true
	} else {
		return false
	}
}

func (w *World) GetSmallTick() int {
	var small big.Int
	small.Set(w.Tick)
	return int(small.Mod(&small, big.NewInt(2<<16)).Uint64())
}

// decreases energy if bacteria is near an acid spot
func (w *World) ApplyAcid(population *evo.Population) {
	for _, b := range population.Bacts {
		if !b.Born {
			continue
		}
		resist := b.GetAcidResist()
		damage := 0.0
		for _, a := range w.Acid {
			dist := dist(a.X, a.Y, b.X, b.Y)
			damage += (a.Con - a.Con*resist) / (dist + 0.001)
		}
		b.Energy = math.Max(0, b.Energy-damage)
	}
}

// returns nearest food
func (w *World) GetNearestFood(b *evo.Bacteria) *Food {
	min_dist := math.Inf(0)
	var result *Food
	for _, f := range w.Food {
		if f.Eaten {
			continue
		}
		dist := dist(b.X, b.Y, f.X, f.Y)
		if dist < min_dist {
			min_dist = dist
			result = f
		}
	}
	return result
}

// returns nearest acid point
func (w *World) GetNearestAcid(b *evo.Bacteria) *Acid {
	min_dist := math.Inf(0)
	var result *Acid
	for _, a := range w.Acid {
		dist := dist(b.X, b.Y, a.X, a.Y)
		if dist < min_dist {
			min_dist = dist
			result = a
		}
	}
	return result
}

// returns the nearest enemy bacteria
func (w *World) GetNearestEnemy(b *evo.Bacteria) *evo.Bacteria {
	min_dist := math.Inf(0)
	var result *evo.Bacteria
	for _, p := range w.Populations {
	outerLoop:
		for _, f := range p.Bacts {
			if f.Chromosome.Author == b.Chromosome.Author {
				continue outerLoop
			}
			dist := dist(b.X, b.Y, f.X, f.Y)
			if dist < min_dist {
				min_dist = dist
				result = f
			}
		}
	}
	return result
}

// returns the nearest fellow bacteria
func (w *World) GetNearestFellow(b *evo.Bacteria) *evo.Bacteria {
	min_dist := math.Inf(0)
	var result *evo.Bacteria
	for _, p := range w.Populations {
	outerLoop:
		for _, f := range p.Bacts {
			if f.Chromosome.Author != b.Chromosome.Author {
				continue outerLoop
			}
			if f == b {
				continue
			}
			if !f.Born {
				continue
			}
			dist := dist(b.X, b.Y, f.X, f.Y)
			if dist < min_dist {
				min_dist = dist
				result = f
			}
		}
	}
	return result
}

// returns the nearest fellow bacteria
func (w *World) GetNearestBact(b *evo.Bacteria) *evo.Bacteria {
	min_dist := math.Inf(0)
	var result *evo.Bacteria
	for _, p := range w.Populations {
		for _, f := range p.Bacts {
			if f == b {
				continue
			}
			dist := dist(b.X, b.Y, f.X, f.Y)
			if dist < min_dist {
				min_dist = dist
				result = f
			}
		}
	}
	return result
}

// returns distance between two points
func dist(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func (w *World) Snapshot() *World {
	braveNewWorld := NewWorld()

	braveNewWorld.Tick = big.NewInt(0)
	braveNewWorld.Tick.Set(w.Tick)

	for _, p := range w.Populations {
		newPop := p.Clone()
		newPop.Chain = nil

		braveNewWorld.Populations = append(braveNewWorld.Populations, newPop)
	}

	for _, f := range w.Food {
		newFood := *f
		braveNewWorld.Food = append(braveNewWorld.Food, &newFood)
	}

	return braveNewWorld
}
