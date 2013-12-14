package engine

import (
	"cryptobact/evo"

	"math"
	"math/rand"
	"math/big"
)

const (
	MAX_ACID_SPOTS = 5
	MAX_ACID_CON   = 0.5
	MAX_CLOT_SPOTS = 5
	MAX_CLOT_DENS  = 0.5
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
	return &World{Tick: big.NewInt(0)}
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
	if w.Notch(w.FoodTicks) {
		return
	}

	for i := 0; i < w.FoodPerTick; i++ {
		x := rand.Float64() * (float64(w.Width) - 1)
		y := rand.Float64() * (float64(w.Height) - 1)
		w.Food = append(w.Food, &Food{x, y, false})
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
		if !b.Born {
			continue
		}
		b.TTL -= int(population.GetGene(b, 17)/10.0 + 1)
	}
}

func (w *World) Step() {
	w.Tick = w.Tick.Add(w.Tick, big.NewInt(1))
}

func (w *World) Notch(notch int) bool {
	if w.Tick.Mod(w.Tick, big.NewInt(int64(notch))) == big.NewInt(0) {
		return true
	} else {
		return false
	}
}

// decreases energy if bacteria is near an acid spot
func (w *World) ApplyAcid(population *evo.Population) {
	for _, b := range population.Bacts {
		if !b.Born {
			continue
		}
		resist := b.GetAcidResist()
		// TODO maybe, delimit by radius
		damage := 0.0
		for _, a := range w.Acid {
			dist := (math.Pow(a.X-b.X, 2.0) + math.Pow(a.Y-b.Y, 2.0))
			damage += (a.Con - a.Con*resist) / (dist + 0.001)
		}
		b.Energy = math.Max(evo.B_MIN_ACID_HEALTH, b.Energy-damage)
	}
}

// descreases speed if bacteria is in clot spot
func (w *World) ApplyClot(population *evo.Population) {
	for _, b := range population.Bacts {
		if !b.Born {
			continue
		}
		resist := b.GetClotResist()
		// TODO maybe, delimit by radius
		slowdown := 0.0
		for _, c := range w.Clot {
			dist := (math.Pow(c.X-b.X, 2.0) + math.Pow(c.Y-b.Y, 2.0))
			slowdown += (c.Density - c.Density*resist) / (dist + 0.001)
		}
		b.Speed = math.Max(evo.B_MIN_SPEED, b.Speed-slowdown)
		b.RotationSpeed = math.Max(evo.B_MIN_ROT_SPEED, b.RotationSpeed-slowdown)
	}
}
