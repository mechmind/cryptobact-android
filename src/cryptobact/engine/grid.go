package engine

import "math"

const(
	MIN_WEIGHT = 0.001
)

type Cell struct {
	Food float64
	Acid float64
	Clot float64
}

type Grid [][]Cell

func (g *Grid) CalcWeights(w *World) {
    for x, row := range *g {
        for y, _ := range row {
			c := Cell{}
			for i := range w.Food {
				c.Food += GetWeight(
					float64(x),
					float64(y),
					w.Food[i].X,
					w.Food[i].Y,
				)
			}
			for i := range w.Acid {
				c.Acid += GetWeight(
					float64(x),
					float64(y),
					w.Acid[i].X,
					w.Acid[i].Y,
				)
			}
			for i := range w.Clot {
				c.Clot += GetWeight(
					float64(x),
					float64(y),
					w.Clot[i].X,
					w.Clot[i].Y,
				)
			}
		}
	}
	return
}

func GetWeight(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	result := 1 / (math.Pow(x2 - x1, 2) + math.Pow(y2 - y1, 2) + 1)
	if result < MIN_WEIGHT {
		return 0
	}
	return result
}
