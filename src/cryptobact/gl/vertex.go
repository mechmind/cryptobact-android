package gl

type Vertex struct {
	X, Y float32
}

type ColoredVertex struct {
	v Vertex
	c Color
}
