package bact

import(
	"cryptobact/engine/grid"
	"cryptobact/engine/world"
)

type Bact struct{
	Owner int
	X float64
	Y float64
	Speed float64
	Food float64
	Acid float64
	Clot float64
	Attack float64
	Fuck float64
	Ttl int
	Energy int
}

type Action interface {
	Apply(w *World)
}

type ActionMove struct {
	Subject Bact
	X float64
	Y float64
}

func (a ActionMove) Apply(world *World) {
	return
}

type ActionAttack struct {
	Subject Bact
	Object Bact
}

func (a ActionAttack) Apply(world *World) {
	return
}

type ActionEat struct {
	Subject Bact
	Object Food
}

func (a ActionEat) Apply(world *World) {
	return
}

type ActionFuck struct {
	Subject Bact
	Object Food
}

func (a ActionFuck) Apply(world *World) {
	return
}

type ActionDie struct {
	Subject Bact
}

func (a ActionDie) Apply(world *World) {
	return
}

func GetAction(grid Grid, world World) Action {
	return Action{}
}
