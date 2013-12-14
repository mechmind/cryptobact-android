package ui

// field with bacterias

import (
	"cryptobact/engine"
	"cryptobact/gl"
)

type Field struct {
	bacs   []*gl.Buffer
	eggs   []*gl.Buffer
	food   *gl.Buffer
	markup *gl.Buffer
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) UpdateBact(cx, cy int, angle float32, color float32)   {}
func (f *Field) UpdateEgg()    {}
func (f *Field) UpdateFood()   {}
func (f *Field) UpdateMarkup() {}

func renderObject(pattern []gl.Vertex, cx, cy float32, other ...float32) []float32 {
	vexs := make([]float32, len(pattern)*(2+len(other)))
	//log.Println("render: coords are", cx, cy)
	step := len(other) + 2
	for idx := 0; idx < len(pattern); idx += step {
		vexs[idx] = cx*STEP + pattern[idx].X
		vexs[idx+1] = cy*STEP + pattern[idx].Y
		copy(vexs[idx+1:], other)
	}
	//log.Println("render: resulting coord set is", vexs)
	return vexs
}
