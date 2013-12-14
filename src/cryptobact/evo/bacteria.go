package evo

import "fmt"

const (
	B_MIN_ACID_HEALTH = 0    // minimum health remaining after acid poisoning
	B_MIN_SPEED       = 0.05 // minimum speed
	B_MIN_ROT_SPEED   = 0.05 // minimum rotation speed
)

type Bacteria struct {
	Chromosome    *Chromosome
	TTL           int
	Energy        float64
	X             float64 `json:"-"`
	Y             float64 `json:"-"`
	Angle         float64 `json:"-"`
	Born          bool
	Speed         float64 `json:"-"`
	RotationSpeed float64 `json:"-"`
	TargetX       float64 `json:"-"`
	TargetY       float64 `json:"-"`
}

func (b *Bacteria) String() string {
	return fmt.Sprintf("{%5.2f; %5.2f} A%3.2f E%6.1f TTL%5d :: %s [%t]\n",
		b.X, b.Y,
		b.Angle,
		b.Energy,
		b.TTL,
		b.Chromosome.DNA,
		b.Born,
	)
}

func (b *Bacteria) GetAggressiveness() float64 {
	// FIXME implement
	return 0.5
}

func (b *Bacteria) GetHunger() float64 {
	// FIXME implement
	return 0.5
}

func (b *Bacteria) GetFertility() float64 {
	// FIXME implement
	return 0.5
}

func (b *Bacteria) GetAcidResist() float64 {
	// FIXME implement
	return 0.5
}

func (b *Bacteria) GetClotResist() float64 {
	// FIXME implement
	return 0.5
}
