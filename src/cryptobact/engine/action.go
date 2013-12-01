package engine

import "log"

import "math"
import "math/rand"
import "cryptobact/evo"

var _ = log.Print

type Action interface {
	Apply()
}

type ActionMove struct {
	X float64
	Y float64
    World *World
    Population *evo.Population
	Bact *evo.Bacteria
}

func (a ActionMove) Apply() {
	b := a.Bact
	alpha := b.Angle
	xt := a.X
	yt := a.Y
	x := b.X
	y := b.Y
	xz := (math.Cos(alpha)) + x
	yz := (math.Sin(alpha)) + y

	ta := math.Sqrt(math.Pow((xz - xt), 2) + math.Pow((yz - yt), 2))
	tb := math.Sqrt(math.Pow((xz - x), 2) + math.Pow((yz - y), 2))
	tc := math.Sqrt(math.Pow((xt - x), 2) + math.Pow((yt - y), 2))

	cos := ((math.Pow(tb, 2) + math.Pow(tc, 2) - math.Pow(ta, 2)) / (2 * tb * tc))
	gamma := math.Acos(cos) * 180 / math.Pi

	direction := "ccw"
	if (math.Cos(xt) - math.Sin(yt)) < 0 {
		direction = "cw"
	}

	if gamma < b.RotationSpeed {
		b.Angle = gamma
		return
	}

	if direction == "cw" {
		b.Angle -= b.RotationSpeed

		if b.Angle < 0 {
			b.Angle += 360
		}
	} else {
		b.Angle += b.RotationSpeed
		if b.Angle > 359 {
			b.Angle -= 360
		}
	}

    dx := (x - xt) / tc + a.Population.GetGene(a.Bact, 6) / 100.0
    dy := (y - yt) / tc + a.Population.GetGene(a.Bact, 7) / 100.0

    b.X += rand.NormFloat64() * 0.1 + dx
    b.Y += rand.NormFloat64() * 0.1 + dy
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
	b.Energy += float64(a.Population.GetGene(b, 12))
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
	child.RotationSpeed = 10.0 + float64(a.Population.GetGene(child, 4) / 20)

    a.Object.Energy -= b_coeff / b_lust * 80
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
    //
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

	if (rand.Intn(300) == 5) {
		for _, b := range population.GetBacts() {
			if b.Energy > 0 && b.Born {
				return ActionFuck{b, world, population, bact}
			}
		}
	}

	// FIXME replace with real target
	//target_x := float64(world.Width) / 2.0
	//target_y := float64(world.Height) / 2.0
	target_x := 1.0
	target_y := 1.0

	return ActionMove{target_x, target_y, world, population, bact}
}
