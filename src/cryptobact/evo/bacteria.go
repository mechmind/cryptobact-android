package evo

import "fmt"

// base values common for all bacterias
const (
	B_BASE_ENERGY            = 200    // energy points at birth
	B_BASE_TTL               = 4000   // time to live in ticks
	B_BASE_SPEED             = 0.7    // speed in pixels per tick
	B_BASE_ROTATION          = 0.5    // rotation speed in degrees per tick
	B_BASE_METABOLISM        = 0.75   // which part of eaten food becomes an energy
	B_BASE_CLOT_RESIST       = 0.5    // clot resistance {0..1}
	B_BASE_ACID_RESIST       = 0.5    // acid resistance {0..1}
	B_BASE_FERTILITY         = 0.5    // fertility rank {0..1}
	B_BASE_DAMAGE            = 2      // physical damage per tick
	B_BASE_LUST              = 0.5    // love to food {0..1}
	B_BASE_GLUT              = 0.5    // love to fuck {0..1}
	B_BASE_AGGRESSION        = 0.5    // aggression {0..1}
	B_BASE_FUCK_ENERGY       = 150    // energy spent to fuck
	B_BASE_MOVE_ENERGY       = 0.001  // energy spent to move
	B_BASE_PROCR_ENERGY      = 0.0001 // energy wiped while procrastinate (per tick)
	B_BASE_EAT_DIST          = 0.05   // maximum eat distance
	B_BASE_FUCK_DIST         = 1.2    // maximum fuck distance
	B_BASE_ATTACK_DIST       = 2.0    // maximum attack distance
	B_BASE_COLLISION_DIST    = 0.8    // collision detection radius
	B_BASE_COLLISION_SPEED   = 0.1    // speed after colision
	B_BASE_COLLISION_INERTIA = 0.8    // speed after colision
	B_BASE_HYSTERIA          = 0.02   // speed after colision
)

// Inertia vector (generated as a result of collision)
type Inertia struct {
	X float64
	Y float64
}

// values specific for current bacteria sample
type Bacteria struct {
	Chromosome *Chromosome
	TTL        int     // current ttl
	Energy     float64 // current energy
	X          float64 `json:"-"`
	Y          float64 `json:"-"`
	Angle      float64 `json:"-"`
	Born       bool
	Labouring  bool
	Inertia    *Inertia `json:"-"`
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
		false,
		&Inertia{},
	}
}

// FIXME get coeff from DNA
func getSelfTtl(d *DNA) int {
	coeff := d.GetNormGene(2)
	result := B_BASE_TTL + B_BASE_TTL*coeff +
		10*d.GetNormGene(0) -
		10*d.GetNormGene(1) +
		2*d.GetNormGene(2) -
		2*d.GetNormGene(3) +
		d.GetNormGene(4) -
		d.GetNormGene(5) +
		d.GetNormGene(6)
	return int(result)
}

func (b *Bacteria) GetSelfTtl() int {
	return getSelfTtl(b.Chromosome.DNA)
}

// FIXME get coeff from DNA
func getSelfEnergy(d *DNA) float64 {
	coeff := d.GetNormGene(1) - d.GetNormGene(3)/(1+d.GetNormGene(14))
	result := B_BASE_ENERGY + B_BASE_ENERGY*coeff
	return result
}

func (b *Bacteria) GetSelfEnergy() float64 {
	return getSelfEnergy(b.Chromosome.DNA)
}

// FIXME get coeff from DNA
func (b *Bacteria) GetSpeed() float64 {
	coeff := (b.Chromosome.DNA.GetNormGene(0) - b.Chromosome.DNA.GetNormGene(2)) /
		(b.Chromosome.DNA.GetNormGene(11) + 0.01)
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
	coeff := b.Chromosome.DNA.GetNormGene(11)
	result := (B_BASE_LUST + B_BASE_LUST*coeff) / (float64(b.TTL)*0.0001 + 0.01)
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetGlut() float64 {
	coeff := b.Chromosome.DNA.GetNormGene(10)
	result := (B_BASE_GLUT + B_BASE_GLUT*coeff) / (float64(b.Energy)*0.001 + 0.01)
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
	coeff := b.Chromosome.DNA.GetNormGene(9)
	result := B_BASE_ACID_RESIST + B_BASE_ACID_RESIST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetFuckEnergy() float64 {
	coeff := b.Chromosome.DNA.GetNormGene(8)
	result := B_BASE_FUCK_ENERGY + B_BASE_FUCK_ENERGY*coeff
	return result
}

func (b *Bacteria) GetHysteria() float64 {
	coeff := b.Chromosome.DNA.GetNormGene(9)
	result := B_BASE_HYSTERIA + B_BASE_HYSTERIA*coeff
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
	coeff := b.Chromosome.DNA.GetNormGene(6)
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
	if b.Energy >= energy_required*(1+b.Chromosome.DNA.GetNormGene(1)) {
		return true
	}
	return false
}

func (b *Bacteria) RenewTTL() {
	b.TTL = int(float64(getSelfTtl(b.Chromosome.DNA)) +
		10000*b.Chromosome.DNA.GetNormGene(10))
}

func (b *Bacteria) String() string {
	form := ""
	if !b.Born {
		form = "EGG"
	}

	return fmt.Sprintf("{%5.2f;%5.2f}A%-5.1fE%-6.1fTTL%-5d :: %s %s",
		b.X, b.Y,
		b.Angle,
		b.Energy,
		b.TTL,
		b.Chromosome.DNA,
		form,
	)
}

// FIXME get coeff from DNA
func (b *Bacteria) GetProcrEnergy() float64 {
	coeff := 0.0
	result := B_BASE_PROCR_ENERGY + B_BASE_PROCR_ENERGY*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetCollisionDist() float64 {
	coeff := 0.0
	result := B_BASE_COLLISION_DIST + B_BASE_COLLISION_DIST*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetCollisionSpeed() float64 {
	coeff := 0.0
	if !b.Born {
		coeff = -0.6
	}
	result := B_BASE_COLLISION_SPEED + B_BASE_COLLISION_SPEED*coeff
	return result
}

func (b *Bacteria) GetCollisionInertia() float64 {
	coeff := 0.0
	if !b.Born {
		coeff = -0.6
	}
	result := B_BASE_COLLISION_INERTIA + B_BASE_COLLISION_INERTIA*coeff
	return result
}

// FIXME get coeff from DNA
func (b *Bacteria) GetMoveEnergy() float64 {
	coeff := b.Chromosome.DNA.GetNormGene(5)
	result := (B_BASE_MOVE_ENERGY + B_BASE_MOVE_ENERGY*coeff) *
		b.Chromosome.DNA.GetNormGene(4)
	return result
}

func (b *Bacteria) GetColor() []byte {
	return []byte{
		byte(int(255.0 * (b.Chromosome.DNA.GetNormGene(0) + b.Chromosome.DNA.GetNormGene(3))) % 255),
		byte(int(255.0 * (b.Chromosome.DNA.GetNormGene(1) + b.Chromosome.DNA.GetNormGene(4))) % 255),
		byte(int(255.0 * (b.Chromosome.DNA.GetNormGene(2) + b.Chromosome.DNA.GetNormGene(5))) % 255),
	}
}

func (b *Bacteria) Clone() *Bacteria {
	newBacteria := *b
	// Chromosome intentionaly left same

	return &newBacteria
}
