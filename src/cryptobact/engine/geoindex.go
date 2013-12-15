package engine

import (
	"log"
	"math"
)

var _ = log.Println

type GeoIndex struct {
	data  map[int]map[int][]Locater
	scale float64
}

type Locater interface {
	Locate() (float64, float64)
}

func NewGeoIndex(scale float64) *GeoIndex {
	return &GeoIndex{
		data:  make(map[int]map[int][]Locater),
		scale: scale,
	}
}

func (g *GeoIndex) Populate(list *[]Locater) {
	for _, l := range *list {
		g.Insert(l)
	}
}

func (g *GeoIndex) Insert(l Locater) {
	x, y := l.Locate()

	if g.data[int(x/g.scale)] == nil {
		g.data[int(x/g.scale)] = make(map[int][]Locater)
	}

	a := g.data[int(x/g.scale)][int(y/g.scale)]
	if a == nil {
		g.data[int(x/g.scale)][int(y/g.scale)] = make([]Locater, 0)
	}
	g.data[int(x/g.scale)][int(y/g.scale)] = append(a, l)
}

func (g *GeoIndex) Quadrant(x float64, y float64) []Locater {
	return g.data[int(x/g.scale)][int(y/g.scale)]
}

func (g *GeoIndex) GetNearest(l Locater) (float64, Locater) {
	x, y := l.Locate()
	min_dist := math.Inf(0)
	var result Locater = nil
	for _, l2 := range g.Quadrant(x, y) {
		x2, y2 := l2.Locate()
		if l == l2 {
			continue
		}
		dist := dist(x, x2, y, y2)
		if dist < min_dist {
			min_dist = dist
			result = l2
		}
	}

	return min_dist, result
}
