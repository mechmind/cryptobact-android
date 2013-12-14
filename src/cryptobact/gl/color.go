package gl

type Color struct {
	R, G, B, A float32
}

func (c Color) RGB() []float32 {
	return []float32{c.R, c.G, c.B}
}

func (c Color) Packed() float32 {
	return 0
}
