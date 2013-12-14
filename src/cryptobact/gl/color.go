package gl

import "unsafe"

type Color struct {
	R, G, B, A float32
}

func (c Color) RGB() []float32 {
	return []float32{c.R, c.G, c.B}
}

func (c Color) Packed() float32 {
	return 0
}

func PackColor(color [3]byte) float32 {
	wword := [4]byte{color[0], color[1], color[2], 1}
	colP := (*float32)(unsafe.Pointer(&wword[0]))
	return *colP
}
