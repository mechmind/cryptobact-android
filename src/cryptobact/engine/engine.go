package engine

import(
	/*"log"*/
	"math"
)

const(
	GRID_WIDTH = 20
	GRID_HEIGHT = 20
	TIME_DELTA = 0.001
)

type Cell struct {
	acid float64
	clot float64
	food float64
}

type Grid [GRID_WIDTH][GRID_HEIGHT]Cell

type Modifier struct {
	x float64
	y float64
	t string
}

type Modifiers []Modifier

func GetMap() Grid {
	var g Grid

	var mods []Modifier
	mods = append(mods, Modifier{0.0, 0.0, "acid"})
	mods = append(mods, Modifier{10.0, 10.0, "acid"})
	mods = append(mods, Modifier{19.0, 19.0, "clot"})
	mods = append(mods, Modifier{10.0, 10.0, "food"})

	g.ApplyModifiers(mods)

	return g
}

func GetModifierWeight(
	x_self float64, y_self float64, x_mod float64, y_mod float64) float64 {
	x_diff := x_mod - x_self
	y_diff := y_mod - y_self
	result := 1 / (math.Pow(x_diff, 2) + math.Pow(y_diff, 2) + 1)

	if result < 0.001 {
		return 0
	}

	return result
}

func (g *Grid) ApplyModifiers(mods Modifiers) {
	for x := 0; x < GRID_WIDTH; x++ {
		for y := 0; y < GRID_HEIGHT; y++ {
			c := Cell{0, 0, 0}
			for i := range mods {
				switch mods[i].t {
					case "acid":
						c.acid += GetModifierWeight(
							float64(x), float64(y), mods[i].x, mods[i].y)
					case "clot":
						c.clot += GetModifierWeight(
							float64(x), float64(y), mods[i].x, mods[i].y)
					case "food":
						c.food += GetModifierWeight(
							float64(x), float64(y), mods[i].x, mods[i].y)
				}
			}
			g[x][y] = c
		}
	}
}
