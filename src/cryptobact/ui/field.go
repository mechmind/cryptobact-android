package ui

// field with bacterias

import (
	"cryptobact/gl"
	"log"
	"math"
)

var ticks int

type Field struct {
	buffers []*gl.Buffer

	// moving field
	vxShader, fragShader, prog uint
	// shader props
	position, mvp, offset, color int
	// offset
	offx, offy float32

	// static field
	svxShader, sfragShader, sprog uint
	// shader props
	sposition, smvp, soffset, scolor, bcolor int
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
	f.prog, err = gl.CreateProgram(gridVXShader, gridFragShader)
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
		gl.ErrPanic()
		gl.GlUseProgram(f.sprog)
		gl.ErrPanic()
		//gl.GlVertexAttribPointer(f.sposition, 2, gl.FLOAT, false, 0, 0)
		gl.GlVertexAttribPointer(f.sposition, 2, gl.FLOAT, false, 12, 0)
		gl.ErrPanic()
		gl.GlVertexAttribPointer(f.scolor, 4, gl.UNSIGNED_BYTE, true, 12, 8)
		gl.ErrPanic()
		gl.GlEnableVertexAttribArray(f.position)
		gl.ErrPanic()
		gl.GlEnableVertexAttribArray(f.scolor)
		gl.ErrPanic()
		gl.GlUniformMatrix4fv(f.smvp, 1, false, mvp)
		gl.ErrPanic()
		return nil
	}

	gridBinder := func(b *gl.Buffer) error {
		//log.Println("field: grid binder step 0")
		gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
		gl.ErrPanic()
		//log.Println("field: grid binder step 1")
		gl.GlUseProgram(f.prog)
		gl.ErrPanic()
		//log.Println("field: grid binder step 2")
		gl.GlVertexAttribPointer(f.position, 2, gl.FLOAT, false, 0, 0)
		gl.ErrPanic()
		//log.Println("field: grid binder step 3")
		gl.GlEnableVertexAttribArray(f.position)
		gl.ErrPanic()
		//log.Println("field: grid binder step 4")
		gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
		gl.ErrPanic()
		//log.Println("field: grid binder done")
		return nil
	}
	f.buffers[ID_MARKUP].Binder = gl.ShaderBinder(gridBinder)
	//gl.GlUniformMatrix4fv(f.mvp, 1, false, mvp)
	// set up grid buffer
	f.buffers[ID_MARKUP].Append(makeGridPoints(X_COUNT*STEP, Y_COUNT*STEP, STEP))
	f.buffers[ID_MARKUP].Flush()

	//	verts := makeGridPoints(1.0, 1.0, 0.3)
	//	verts = append(verts, 0.1, 0.1, 0.5, 0.5, 1.0, 1.0, 2.0, 2.0, 10.0, 10.0, 50.0, 50.0)
	//	b := f.buffers[ID_MARKUP]
	//	b.BufLen = len(verts)
	//	gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	//	gl.ErrPanic()
	//	gl.GlBufferData(gl.ARRAY_BUFFER, verts, gl.STATIC_DRAW)
	//	gl.ErrPanic()
	//	gl.GlVertexAttribPointer(f.position, 2, gl.FLOAT, false, 0, 0)
	//	gl.ErrPanic()

	// ***** BACTERIA SHADERS (with own prog) ******
	f.sprog, err = gl.CreateProgram(bactVXShader, bactFragShader)
	if err != nil {
		return err
	}
	log.Println("field: step 1")
	gl.ErrPanic()
	//
	f.sposition, _ = gl.GlGetAttribLocation(f.sprog, "position")
	gl.ErrPanic()
	f.smvp, _ = gl.GlGetUniformLocation(f.sprog, "mvp")
	gl.ErrPanic()
	f.scolor, _ = gl.GlGetAttribLocation(f.sprog, "lcolor")
	gl.ErrPanic()
	f.soffset, _ = gl.GlGetUniformLocation(f.sprog, "offset")
	gl.ErrPanic()
	f.bcolor, _ = gl.GlGetUniformLocation(f.sprog, "color")
	gl.ErrPanic()
	gl.GlUseProgram(f.sprog)
	log.Println("field: prog and shader inputs:", f.sprog, f.sposition, f.smvp, f.scolor, f.soffset, f.bcolor)
	log.Println("field: step 2")
	//
	gl.GlEnableVertexAttribArray(f.sposition)
	gl.ErrPanic()
	gl.GlEnableVertexAttribArray(f.scolor)
	gl.ErrPanic()
	log.Println("field: step 3")
	//	sgridBinder := func(b *gl.Buffer) error {
	//		gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	//		gl.ErrPanic()
	//		gl.GlVertexAttribPointer(f.sposition, 2, gl.FLOAT, false, 0, 0)
	//		gl.ErrPanic()
	//		gl.GlEnableVertexAttribArray(f.sposition)
	//		gl.ErrPanic()
	//		gl.GlUniformMatrix4fv(f.smvp, 1, false, mvp)
	//		gl.ErrPanic()
	//		gl.GlUseProgram(f.sprog)
	//		gl.ErrPanic()
	//		return nil
	//	}
	f.buffers[ID_BACTERIA_BODY].Binder = gl.ShaderBinder(bactBinder)
	f.buffers[ID_BACTERIA_EYES].Binder = gl.ShaderBinder(bactBinder)

	// temporary bacteria body
	//	gl.GlUniformMatrix4fv(f.smvp, 1, false, mvp)
	return nil
}

func (f *Field) Draw(mvp []float32) {
	f.buffers[ID_MARKUP].Binder(f.buffers[ID_MARKUP])
	//log.Println("field: rendering grid", f.buffers[ID_MARKUP].BufLen)

	// remove after debug
	b := f.buffers[ID_MARKUP]
	gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	gl.ErrPanic()
	gl.GlUniform3f(f.color, 1.0, 1.0, 1.0)
	gl.GlUniform2f(f.offset, f.offx, f.offy)
	//gl.GlDrawArrays(gl.POINTS, 0, b.BufLen)
	//gl.GlClear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.GlDrawArrays(gl.POINTS, 0, b.BufLen)

	//f.buffers[ID_BACTERIA_BODY].Binder(f.buffers[ID_BACTERIA_BODY])
	//gl.DrawArrays

	// ********* bacterias **********
	//	//log.Println("field: rendering grid", f.buffers[ID_MARKUP].BufLen)
	//
	//	// remove after debug

	b = f.buffers[ID_BACTERIA_BODY]
	gl.GlBindBuffer(gl.ARRAY_BUFFER, b.GlBuffer)
	f.buffers[ID_BACTERIA_BODY].Binder(f.buffers[ID_BACTERIA_BODY])

	gl.GlUniform2f(f.soffset, f.offx, f.offy)
	gl.GlDrawArrays(gl.TRIANGLES, 0, b.BufLen)
}

func (f *Field) FlushAll() {
	for _, buf := range f.buffers {
		buf.Flush()
	}
}

func (f *Field) UpdateBact(cx, cy float32, angle float32, color [3]byte) {
	// update body
	bactbuf := renderRotatedObject(mainSet[ID_BACTERIA_BODY].verts,
		cx, cy, angle, gl.PackColor(color))
	//bactbuf = renderObject(mainSet[ID_BACTERIA_BODY].verts, fticks, 10)
	eyebuf := renderRotatedObject(mainSet[ID_BACTERIA_EYES].verts,
		cx, cy, angle, gl.PackColor([3]byte{0, 0, 0}))
	f.buffers[ID_BACTERIA_BODY].Append(eyebuf)

	eyebuf = renderRotatedObject(mainSet[ID_BACTERIA_EYE_SHINE].verts,
		cx, cy, angle, gl.PackColor([3]byte{255, 255, 255}))
	f.buffers[ID_BACTERIA_BODY].Append(eyebuf)

	eyebuf = renderRotatedObject(mainSet[ID_BACTERIA_EYE_STARK].verts,
		cx, cy, angle, gl.PackColor([3]byte{0, 0, 255}))
	f.buffers[ID_BACTERIA_BODY].Append(eyebuf)
	f.buffers[ID_BACTERIA_BODY].Append(bactbuf)
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
	vexs := make([]float32, len(pattern)+len(pattern)/2*len(other))
	//log.Println("render: coords are", cx, cy)
	step := len(other) + 2
	var pidx int
	for idx := 0; idx < len(vexs); idx += step {
		vexs[idx] = cx*STEP + pattern[pidx]*3
		vexs[idx+1] = cy*STEP + pattern[pidx+1]*3
		copy(vexs[idx+2:], other)
		pidx += 2
	}
	//log.Println("render: resulting coord set is", vexs)
	return vexs
}

func renderRotatedObject(pattern []float32, cx, cy, angle float32, other ...float32) []float32 {
	vexs := make([]float32, len(pattern)+len(pattern)/2*len(other))
	//log.Println("render: coords are", cx, cy)
	step := len(other) + 2
	var pidx int
	angle = (angle + 90) * math.Pi / 180
	sin := float32(math.Sin(float64(angle)))
	cos := float32(math.Cos(float64(angle)))
	for idx := 0; idx < len(vexs); idx += step {
		lx := pattern[pidx] * 3
		ly := pattern[pidx+1] * 3
		vexs[idx] = cx*STEP + lx*cos - ly*sin
		vexs[idx+1] = cy*STEP + lx*sin + ly*cos
		copy(vexs[idx+2:], other)
		pidx += 2
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
