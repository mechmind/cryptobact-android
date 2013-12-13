package main

import (
	"cryptobact/engine"
	"fmt"
	"log"
)

type Updater struct{}

func (f Updater) Update(w *engine.World) {
	maxCol := 40
	maxRow := 40

	startX := -10
	endX := 10

	startY := -10
	endY := 10

	stepX := float64(endX-startX) / float64(maxCol)
	stepY := float64(endY-startY) / float64(maxRow)

	fmt.Printf("\n\nTICK: %d\n", w.Tick)
	fmt.Printf("RANGE: X[%d, %d] Y[%d, %d]\n", startX, endX, startY, endY)

	for y := 0; y < maxRow; y++ {
		for x := 0; x < maxCol; x++ {
			printed := false
			for _, p := range w.Populations {
				if printed {
					break
				}
				for _, b := range p.Bacts {
					if !between(b.X, b.Y, stepX, stepY, x, y) {
						continue
					}

					if b.Born {
						//     70
						//    7  0
						//   6    1
						//    5  2
						//     43
						if b.Angle < 45 {
							fmt.Print("0")
						} else if b.Angle < 90 {
							fmt.Print("1")
						} else if b.Angle < 135 {
							fmt.Print("2")
						} else if b.Angle < 180 {
							fmt.Print("3")
						} else if b.Angle < 225 {
							fmt.Print("4")
						} else if b.Angle < 270 {
							fmt.Print("5")
						} else if b.Angle < 315 {
							fmt.Print("6")
						} else {
							fmt.Print("7")
						}
					} else {
						fmt.Print("E")
					}
					printed = true
					break
				}
			}

			if printed {
				continue
			}

			for _, f := range w.Food {
				if !between(f.X, f.Y, stepX, stepY, x, y) {
					continue
				}

				fmt.Print("F")
				printed = true
				break
			}

			if !printed {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func between(x float64, y float64, stepX float64, stepY float64, i int, j int) bool {
	if float64(j-1)*stepY >= y {
		return false
	}

	if float64(j)*stepY <= y {
		return false
	}

	if float64(i-1)*stepX >= x {
		return false
	}

	if float64(i)*stepX <= x {
		return false
	}

	return true
}

func main() {
	log.Println("testing engine")
	u := Updater{}
	engine.Loop(u)
}
