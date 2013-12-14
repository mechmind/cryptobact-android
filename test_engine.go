package main

import (
	"cryptobact/engine"

	"fmt"
	"log"
	"os"
)

type Updater struct{}

func (f Updater) Update(w *engine.World) {
	DrawMap(w.Snapshot())
}

func DrawMap(w *engine.World) {
	return
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
						fmt.Printf("%c", "01234566"[int(b.Angle)/45])
					} else {
						fmt.Print("o")
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

			for _, a := range w.Acid {
				if !between(a.X, a.Y, stepX, stepY, x, y) {
					continue
				}

				fmt.Print("A")
				printed = true
				break
			}

			for _, c := range w.Clot {
				if !between(c.X, c.Y, stepX, stepY, x, y) {
					continue
				}

				fmt.Print("C")
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
	f, _ := os.OpenFile("bact.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)

	log.Println("testing engine")
	u := Updater{}
	engine.Loop(u)
}
