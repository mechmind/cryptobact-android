package bact

import "math/rand"
import . "cryptobact/evo/bacteria"
import . "cryptobact/engine/grid"
import . "cryptobact/engine/world"

type Action interface {
	Apply(b *Bact, w *World)
}

type ActionMove struct {
	X float64
	Y float64
	Angle int
}

func (a ActionMove) Apply(b *Bact, w *World) {
	b.X = a.X
	b.Y = a.Y
	b.Angle = a.Angle
}

type ActionAttack struct {
	Object *Bact
	Damage int
}

func (a ActionAttack) Apply(b *Bact, w *World) {
	a.Object.Energy -= a.Damage
}

type ActionEat struct {
	Object *Food
}

func (a ActionEat) Apply(b *Bact, w *World) {
	a.Object.Eaten = true
}

type ActionFuck struct {
	Object *Bact
}

func (a ActionFuck) Apply(b *Bact, w *World) {
	// FIXME implement
	//child := w.MyPopulation.Fuck(b, a.Object)
	//a_nrg := w.MyPopulation.GetGenome(a.Object)
	//b_nrg := w.MyPopulation.GetGenome(b)
	return
}

type ActionDie struct {
}

func (a ActionDie) Apply(b *Bact, w *World) {
	w.MyPopulation.Kill(b)
}

func GetAction(bact *Bacteria, grid *Grid, world *World) Action {
	// FIXME rewrite without random
	Action.Bact = bact
	actions := []string{"move", "attack", "eat", "fuck", "die"}

	if bact.TTL == 0 {
		return ActionDie{}
	}

	if (rand.Intn(10) == 5) {
		for k, b := range world.MyPopulation {
			if b.Energy > 0 {
				return ActionAttack{b, 30}
			}
		}
	}

	if (rand.Intn(10) == 5) {
		for k, f := range world.Food {
			if f.Eaten == false {
				return ActionEat(f)
			}
		}
	}

	x := (rand.Float64() * 2) - 1
	y := (rand.Float64() * 2) - 1
	for ((x + bact.X) < 0 || (y + bact.Y) < 0 || (x + bact.X) > world.Width ||
		(y + bact.Y) > world.Height) {
		x := (rand.Float64() * 2) - 1
		y := (rand.Float64() * 2) - 1
	}

	return ActionMove{x, y, 0}
}
