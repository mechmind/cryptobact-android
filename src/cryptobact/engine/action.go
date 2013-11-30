package engine

import "math/rand"
import "cryptobact/evo"

type Action interface {
	Apply(b *evo.Bacteria, w *World)
}

type ActionMove struct {
	X float64
	Y float64
	Angle float64
}

func (a ActionMove) Apply(b *evo.Bacteria, w *World) {
	b.X = a.X
	b.Y = a.Y
	b.Angle = a.Angle
}

type ActionAttack struct {
	Object *evo.Bacteria
	Damage float64
}

func (a ActionAttack) Apply(b *evo.Bacteria, w *World) {
	a.Object.Energy -= a.Damage
}

type ActionEat struct {
	Object *Food
}

func (a ActionEat) Apply(b *evo.Bacteria, w *World) {
	a.Object.Eaten = true
}

type ActionFuck struct {
	Object *evo.Bacteria
}

func (a ActionFuck) Apply(b *evo.Bacteria, w *World) {
	// FIXME implement
    child := w.MyPopulation.Fuck(b, a.Object)
    a_coeff := float64(w.MyPopulation.GetGene(a.Object, 0))
    b_coeff := float64(w.MyPopulation.GetGene(b, 0))
    a_lust := float64(w.MyPopulation.GetAttitude(a.Object, "lust"))
    b_lust := float64(w.MyPopulation.GetAttitude(b, "lust"))

    child.X = (a.Object.X + b.X) / 2
    child.Y = (a.Object.Y + b.Y) / 2

    a.Object.Energy -= b_coeff / b_lust * 10
    b.Energy -= a_coeff / a_lust * 10
}

type ActionDie struct {}

func (a ActionDie) Apply(b *evo.Bacteria, w *World) {
	w.MyPopulation.Kill(b)
}

func GetAction(bact *evo.Bacteria, grid *Grid, world *World) Action {
	// FIXME rewrite without random
	//actions := []string{"move", "attack", "eat", "fuck", "die"}

	if bact.TTL <= 0 {
		return ActionDie{}
	}

	if (rand.Intn(10) == 5) {
		for _, b := range world.MyPopulation.GetBacts() {
			if b.Energy > 0 {
				return ActionAttack{b, 30}
			}
		}
	}

	if (rand.Intn(10) == 5) {
		for _, f := range world.Food {
			if f.Eaten == false {
				return ActionEat{f}
			}
		}
	}

	x := (rand.Float64() * 2) - 1
	y := (rand.Float64() * 2) - 1
	for ((x + bact.X) < 0 || (y + bact.Y) < 0 || (x + bact.X) > float64(world.Width) ||
		(y + bact.Y) > float64(world.Height)) {
        x = (rand.Float64() * 2) - 1
        y = (rand.Float64() * 2) - 1
	}

	return ActionMove{x + bact.X, y + bact.Y, 0}
}
