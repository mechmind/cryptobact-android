package evo

import "fmt"

// base values common for all bacterias
const (
	B_BASE_ENERGY      = 10   // energy points at birth
	B_BASE_TTL         = 6000 // time to live in ticks
	B_BASE_SPEED       = 0.05 // speed in pixels per tick
	B_BASE_ROTATION    = 0.05 // rotation speed in degrees per tick
	B_BASE_METABOLISM  = 0.75 // which part of eaten food becomes an energy
	B_BASE_CLOT_RESIST = 0.5  // clot resistance {0..1}
	B_BASE_ACID_RESIST = 0.5  // acid resistance {0..1}
	B_BASE_FERTILITY   = 0.5  // fertility rank {0..1}
	B_BASE_DAMAGE      = 2    // physical damage per tick
	B_BASE_LUST        = 0.5  // love to food {0..1}
	B_BASE_GLUT        = 0.5  // love to fuck {0..1}
	B_BASE_AGGRESSION  = 0.5  // aggression {0..1}
	B_BASE_FUCK_ENERGY = 100  // energy points required to fuck
	B_BASE_EAT_DIST    = 5    // maximum eat distance
	B_BASE_FUCK_DIST   = 5    // maximum fuck distance
	B_BASE_ATTACK_DIST = 5    // maximum attack distance
)

// values specific for current bacteria sample
type Bacteria struct {
	Chromosome    *Chromosome
	TTL           int     // current ttl
	Energy        float64 // current energy
	X             float64 `json:"-"`
	Y             float64 `json:"-"`
	Angle         float64 `json:"-"`
	Born          bool
	Speed         float64 `json:"-"`
	RotationSpeed float64 `json:"-"`
	TargetX       float64 `json:"-"`
	TargetY       float64 `json:"-"`
}

func NewBacteria(c *Chromosome) *Bacteria {
	ttl := getSelfTtl(c.DNA)
	energy := getSelfEnergy(c.DNA)

	return &Bacteria{
		c,
		ttl,
		energy,
		0.0,
		0.0,
		0.0,
		false,
		0.0,
		0.0,
		0.0,
		0.0,
	}
}

// FIXME get coeff from DNA
func getSelfTtl(d *DNA) int {
	coeff := 0
	result := B_BASE_TTL + B_BASE_TTL*coeff
	return result
}

func (b *Bacteria) GetSelfTtl() int {
	return getSelfTtl(b.Chromosome.DNA)
}

// FIXME get coeff from DNA
func getSelfEnergy(d *DNA) float64 {
	coeff := 0.0
	result := B_BASE_ENERGY + B_BASE_ENERGY*coeff
	return result
}

func (b *Bacteria) GetSelfEnergy() float64 {
	return getSelfEnergy(b.Chromosome.DNA)
}

// FIXME get coeff from DNA
func (b *Bacteria) GetSpeed() float64 {
	coeff := 0.0
	result := B_BASE_SPEED + B_BASE_SPEED*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetRotation() float64 {
	coeff := 0.0
	result := B_BASE_ROTATION + B_BASE_ROTATION*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetDamage() float64 {
	coeff := 0.0
	result := B_BASE_DAMAGE + B_BASE_DAMAGE*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetLust() float64 {
	coeff := 0.0
	result := B_BASE_LUST + B_BASE_LUST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetGlut() float64 {
	coeff := 0.0
	result := B_BASE_GLUT + B_BASE_GLUT*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetAggression() float64 {
	coeff := 0.0
	result := B_BASE_AGGRESSION + B_BASE_AGGRESSION*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetAcidResist() float64 {
	coeff := 0.0
	result := B_BASE_ACID_RESIST + B_BASE_ACID_RESIST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetFuckEnergy() float64 {
	coeff := 0.0
	result := B_BASE_FUCK_ENERGY + B_BASE_FUCK_ENERGY*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetMetabolism() float64 {
	coeff := 0.0
	result := B_BASE_METABOLISM + B_BASE_METABOLISM*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetEatDist() float64 {
	coeff := 0.0
	result := B_BASE_EAT_DIST + B_BASE_EAT_DIST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetFuckDist() float64 {
	coeff := 0.0
	result := B_BASE_FUCK_DIST + B_BASE_FUCK_DIST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetAttackDist() float64 {
	coeff := 0.0
	result := B_BASE_ATTACK_DIST + B_BASE_ATTACK_DIST*coeff
	return result
}

func (b *Bacteria) CanFuck() bool {
	energy_required := b.GetFuckEnergy()
	if b.Energy >= energy_required {
		return true
	}
	return false
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

func (b *Bacteria) Clone() *Bacteria {
	newBacteria := *b
	// Chromosome intentionaly left same

	return &newBacteria
}
