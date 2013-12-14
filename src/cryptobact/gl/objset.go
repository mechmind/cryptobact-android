package gl

type shaderBinder func(set *Buffer) error

type Buffer struct {
	glType   uint
	glBuffer uint
	buf      []float32
	binder   shaderBinder
}

func NewBuffer(glType uint, binder shaderBinder) *Buffer {
	glBuf, _ := GlGenBuffer() // FIXME: handle error
	return &Buffer{glType, glBuf, nil, binder}
}

func (b *Buffer) Append(data []float32) {
	b.buf = append(b.buf, data...)
}

// flush data to opengl
// TODO: implement incremental update
func (b *Buffer) Flush() error {
	// upload vertice + meta to opengl
	GlBindBuffer(ARRAY_BUFFER, b.glBuffer)
	GlBufferData(b.glType, b.buf, STATIC_DRAW)
	// FIXME: rebind shader attrs
	b.buf = b.buf[:0]
	return nil
}
