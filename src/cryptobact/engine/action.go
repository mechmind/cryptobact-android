package engine

import "math/rand"
import "cryptobact/evo"

type Action interface {
	Apply()
}

type ActionMove struct {
	X float64
	Y float64
	Angle float64
    World *World
    Population *evo.Population
    Bact *evo.Bacteria
}

func (a ActionMove) Apply() {
    b := a.Bact
	b.X = a.X
	b.Y = a.Y
	b.Angle = a.Angle
}

type ActionAttack struct {
	Object *evo.Bacteria
	Damage float64
    World *World
    Population *evo.Population
    Bact *evo.Bacteria
}

func (a ActionAttack) Apply() {
	a.Object.Energy -= a.Damage
}

type ActionEat struct {
	Object *Food
    World *World
    Population *evo.Population
    Bact *evo.Bacteria
}

func (a ActionEat) Apply() {
    b := a.Bact
	b.Energy += float64(a.Population.GetGene(b, 12)) * 10.0
	a.Object.Eaten = true
}

type ActionFuck struct {
	Object *evo.Bacteria
    World *World
    Population *evo.Population
    Bact *evo.Bacteria
}

func (a ActionFuck) Apply() {
    b := a.Bact
    child := a.Population.Fuck(b, a.Object)
    a_coeff := float64(a.Population.GetGene(a.Object, 0))
    b_coeff := float64(a.Population.GetGene(b, 0))
    a_lust := float64(a.Population.GetAttitude(a.Object, "lust"))
    b_lust := float64(a.Population.GetAttitude(b, "lust"))

    child.X = (a.Object.X + b.X) / 2
    child.Y = (a.Object.Y + b.Y) / 2
	child.TTL = int(10000 * float64(a.Population.GetGene(child, 7)) / 10)
	child.Energy = 1000 * float64(a.Population.GetGene(child, 11)) / 10

    a.Object.Energy -= b_coeff / b_lust * 4
    b.Energy -= a_coeff / a_lust * 4
}

type ActionDie struct {
    World *World
    Population *evo.Population
    Bact *evo.Bacteria
}

func (a ActionDie) Apply() {
    b := a.Bact
	a.Population.Kill(b)
}

func GetAction(population *evo.Population, bact *evo.Bacteria, grid *Grid,
        world *World) Action {
	// FIXME rewrite without random
	//actions := []string{"move", "attack", "eat", "fuck", "die"}

	if bact.TTL <= 0 || bact.Energy < 0 {
		return ActionDie{world, population, bact}
	}

	if (rand.Intn(10) == 5) {
		for _, b := range population.GetBacts() {
			if b.Energy > 0 && b.Born {
				return ActionAttack{b, 30, world, population, bact}
			}
		}
	}

	if (rand.Intn(10) == 5) {
		for _, f := range world.Food {
			if f.Eaten == false {
				return ActionEat{f, world, population, bact}
			}
		}
	}

	if (rand.Intn(30) == 5) {
		for _, b := range population.GetBacts() {
			if b.Energy > 0 && b.Born {
				return ActionFuck{b, world, population, bact}
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

	return ActionMove{x + bact.X, y + bact.Y, 0, world, population, bact}
}
