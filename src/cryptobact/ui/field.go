package ui

// field with bacterias

import (
	"cryptobact/gl"
	"log"
	"math"
)

type Field struct {
	buffers                    []*gl.Buffer
	vxShader, fragShader, prog uint
	// shader props
	position, mvp, offset, color int
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) Init(mvp []float32) error {
	f.buffers = make([]*gl.Buffer, TOTAL_IDS)
	for id := range f.buffers {
		f.buffers[id] = gl.NewBuffer(mainSet[id].glType, nil)
	}

	var err error
	// compile and link shaders
	f.prog, err = gl.CreateProgram(bacteriaVXShader, bacteriaFragShader)
	gl.ErrPanic()
	if err != nil {
		return err
	}
	f.position, _ = gl.GlGetAttribLocation(f.prog, "position")
	gl.ErrPanic()
	f.offset, _ = gl.GlGetUniformLocation(f.prog, "offset")
	gl.ErrPanic()
	f.mvp, _ = gl.GlGetUniformLocation(f.prog, "mvp")
	gl.ErrPanic()
	f.color, _ = gl.GlGetUniformLocation(f.prog, "color")
	gl.GlUseProgram(f.prog)
	log.Println("field: prog and shader inputs:", f.prog, f.position, f.offset, f.mvp, f.color)

	gl.GlEnableVertexAttribArray(f.position)
	gl.ErrPanic()
	// transformation matrix
	bactBinder := func(b *gl.Buffer) error {
		gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
		gl.GlVertexAttribPointer(f.position, 2, gl.FLOAT, false, 0, 0)
		gl.GlEnableVertexAttribArray(f.position)
		gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
		gl.GlUseProgram(f.prog)
		gl.ErrPanic()
		return nil
	}

	gridBinder := func(b *gl.Buffer) error {
		gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
		gl.ErrPanic()
		gl.GlVertexAttribPointer(f.position, 2, gl.FLOAT, false, 0, 0)
		gl.ErrPanic()
		gl.GlEnableVertexAttribArray(f.position)
		gl.ErrPanic()
		gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
		gl.ErrPanic()
		gl.GlUseProgram(f.prog)
		gl.ErrPanic()
		return nil
	}
	f.buffers[ID_BACTERIA_BODY].Binder = gl.ShaderBinder(bactBinder)
	f.buffers[ID_MARKUP].Binder = gl.ShaderBinder(gridBinder)
	//gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
	// set up grid buffer
	//f.buffers[ID_MARKUP].Append(makeGridPoints(X_COUNT*STEP, Y_COUNT*STEP, STEP))
	//f.buffers[ID_MARKUP].Flush()

	verts := makeGridPoints(1.0, 1.0, 0.05)
	verts = append(verts, 0.1, 0.1, 0.5, 0.5, 1.0, 1.0, 2.0, 2.0, 10.0, 10.0, 50.0, 50.0)
	b := f.buffers[ID_MARKUP]
	b.BufLen = len(verts)
	gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	gl.ErrPanic()
	gl.GlBufferData(gl.ARRAY_BUFFER, verts, gl.STATIC_DRAW)
	gl.ErrPanic()
	return nil
}

func (f *Field) Draw(mvp []float32) {
	//f.buffers[ID_MARKUP].Binder(f.buffers[ID_MARKUP])
	log.Println("field: rendering grid", f.buffers[ID_MARKUP].BufLen)

	// remove after debug
	b := f.buffers[ID_MARKUP]
	gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	gl.ErrPanic()
	gl.GlEnableVertexAttribArray(f.position)
	gl.ErrPanic()
	gl.GlVertexAttribPointer(f.position, 2, gl.FLOAT, false, 0, 0)
	gl.ErrPanic()
	//gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
	//gl.ErrPanic()
	//gl.GlUseProgram(f.prog)
	//gl.ErrPanic()
	gl.GlUniform3f(f.color, 1.0, 1.0, 1.0)
	//gl.GlDrawArrays(gl.POINTS, 0, b.BufLen)
	gl.GlDrawArrays(gl.TRIANGLES, 0, b.BufLen)

	//f.buffers[ID_BACTERIA_BODY].Binder(f.buffers[ID_BACTERIA_BODY])
	//gl.DrawArrays
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
