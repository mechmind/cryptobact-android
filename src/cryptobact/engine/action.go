package engine

import (
	"cryptobact/evo"
	"math"
	"math/rand"
)

type Decision struct {
	Weight float64
}

type Action interface {
	Apply()
}

type ActionProcrastinate struct {
	Bact *evo.Bacteria
}

func (a ActionProcrastinate) Apply() {
	a.Bact.Energy -= a.Bact.GetProcrEnergy()
}

type ActionMove struct {
	Bact *evo.Bacteria
	X    float64
	Y    float64
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

	ta := math.Sqrt(math.Pow((xz-xt), 2) + math.Pow((yz-yt), 2))
	tb := math.Sqrt(math.Pow((xz-x), 2) + math.Pow((yz-y), 2))
	tc := math.Sqrt(math.Pow((xt-x), 2) + math.Pow((yt-y), 2))

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

	dx := (xt - x) / math.Abs(x-xt) * b.GetSpeed() / 100.0
	dy := (yt - y) / math.Abs(y-yt) * b.GetSpeed() / 100.0

	b.X += rand.NormFloat64()*0.01 + dx
	b.Y += rand.NormFloat64()*0.01 + dy
}

type ActionAttack struct {
	Bact   *evo.Bacteria
	Object *evo.Bacteria
}

func (a ActionAttack) Apply() {
	a.Object.Energy -= a.Bact.GetDamage()
}

type ActionEat struct {
	Bact   *evo.Bacteria
	Object *Food
}

func (a ActionEat) Apply() {
	b := a.Bact
	b.Energy += b.GetMetabolism() * a.Object.Calories
	a.Object.Eaten = true
}

type ActionFuck struct {
	Bact   *evo.Bacteria
	Object *evo.Bacteria
	P      *evo.Population
}

func (a ActionFuck) Apply() {
	child := a.P.Fuck(a.Bact, a.Object)
	child.X = (a.Object.X + a.Bact.X) / 2
	child.Y = (a.Object.Y + a.Bact.Y) / 2
	a.Bact.Energy -= a.Bact.GetFuckEnergy()
}

type ActionDie struct {
	Bact *evo.Bacteria
	P    *evo.Population
}

func (a ActionDie) Apply() {
	b := a.Bact
	a.P.Kill(b)
}

func GetAction(p *evo.Population, b *evo.Bacteria, w *World) Action {
	if b.TTL <= 0 || b.Energy <= 0 {
		return ActionDie{b, p}
	}

	// params:
	//   aggressions {0..1}
	//   lust {0..1}
	//   glut {0..1}
	// resists:
	//   acid {0..1}
	//   clot {0..1}

	// get nearest food, fellow, enemy and acid
	// calculate weight for every action
	// perform the action with the gratest weight
	var action Action
	max_weight := 0.0
	if n_food := w.GetNearestFood(b); n_food != nil {
		food_dist := dist(n_food.X, n_food.Y, b.X, b.Y)
		weight := b.GetLust() / food_dist
		max_weight = weight
		if b.GetEatDist() <= food_dist {
			action = ActionEat{b, n_food}
		} else {
			action = ActionMove{b, n_food.X, n_food.Y}
		}
	}

	if n_fellow := w.GetNearestFellow(b); b.CanFuck() && n_fellow != nil {
		fellow_dist := dist(n_fellow.X, n_fellow.Y, b.X, b.Y)
		weight := b.GetGlut() / fellow_dist
		if weight > max_weight {
			max_weight = weight
			if b.GetFuckDist() <= fellow_dist {
				action = ActionFuck{b, n_fellow, p}
			} else {
				action = ActionMove{b, n_fellow.X, n_fellow.Y}
			}
		}
	}

	if n_enemy := w.GetNearestEnemy(b); n_enemy != nil {
		enemy_dist := dist(n_enemy.X, n_enemy.Y, b.X, b.Y)
		weight := b.GetAggression() / enemy_dist
		if weight > max_weight {
			max_weight = weight
			if b.GetAttackDist() <= enemy_dist {
				action = ActionAttack{b, n_enemy}
			} else {
				action = ActionMove{b, n_enemy.X, n_enemy.Y}
			}
		}
	}

	if n_acid := w.GetNearestAcid(b); n_acid != nil {
		acid_dist := dist(n_acid.X, n_acid.Y, b.X, b.Y)
		weight := b.GetAcidResist() / acid_dist
		if weight > max_weight {
			max_weight = weight
			x, y := getRunawayPoint(b, n_acid.X, n_acid.Y)
			action = ActionMove{b, x, y}
		}
	}

	if action == nil {
		action = ActionProcrastinate{b}
	}

	return action
}

func getRunawayPoint(b *evo.Bacteria, x float64, y float64) (float64, float64) {
	// FIXME implement
	return 0.0, 0.0
}
