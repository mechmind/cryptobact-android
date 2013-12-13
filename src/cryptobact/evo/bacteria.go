package evo

import "fmt"

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
