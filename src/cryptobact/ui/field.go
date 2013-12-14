package ui

// field with bacterias

import (
	"cryptobact/gl"
	"math"
)

type Field struct {
	buffers []*gl.Buffer
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) Init() error {
	f.buffers = make([]*gl.Buffer, TOTAL_IDS)
	return nil
}

func (f *Field) UpdateBact(cx, cy float32, angle float32, color [3]byte) {
	colorf := gl.PackColor(color)
	// update body
	data := renderObject(mainSet[ID_BACTERIA_BODY].verts, cx, cy, angle, colorf)
	f.buffers[ID_BACTERIA_BODY].Append(data)
}

func (f *Field) UpdateEgg(cx, cy float32, color [3]byte) {
	colorf := gl.PackColor(color)
	// update body
	data := renderObject(mainSet[ID_EGG].verts, cx, cy, colorf)
	f.buffers[ID_EGG].Append(data)
}

func (f *Field) UpdateFood(cx, cy float32) {
	colorf := gl.PackColor(mainSet[ID_FOOD].color)
	// update body
	data := renderObject(mainSet[ID_FOOD].verts, cx, cy, colorf)
	f.buffers[ID_FOOD].Append(data)
}

func (f *Field) UpdateMarkup(cx, cy float32) {
	colorf := gl.PackColor(mainSet[ID_MARKUP].color)
	// update body
	data := renderObject(mainSet[ID_MARKUP].verts, cx, cy, colorf)
	f.buffers[ID_MARKUP].Append(data)
}

func renderObject(pattern []float32, cx, cy float32, other ...float32) []float32 {
	vexs := make([]float32, len(pattern)*(1+len(other)))
	//log.Println("render: coords are", cx, cy)
	step := len(other) + 2
	for idx := 0; idx < len(pattern); idx += step {
		vexs[idx] = cx*STEP + pattern[idx]
		vexs[idx+1] = cy*STEP + pattern[idx+1]
		copy(vexs[idx+1:], other)
	}
	//log.Println("render: resulting coord set is", vexs)
	return vexs
}

func makeGridPoints(limX, limY, step float32) []float32 {
	data := make([]float32, 0, int(math.Ceil(float64(limX)*float64(limY)/(float64(step)*float64(step))+4)*4))

	var nextX, nextY float32
	for nextX = 0.0; nextX < limX+0.1; nextX += step {
		for nextY = 0.0; nextY < limY+0.1; nextY += step {
			data = append(data, nextX, nextY)
			//            data = append(data, nextX - CROSS_HALFSIZE, nextY, nextX + CROSS_HALFSIZE, nextY)
			//            data = append(data, nextX, nextY - CROSS_HALFSIZE, nextX, nextY + CROSS_HALFSIZE)
		}
	}
	return data
}
