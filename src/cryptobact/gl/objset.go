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
	return &Buffer{nil, glType, glBuf, nil, binder}
}

// flush data to opengl
// TODO: implement incremental update
func (b *Buffer) Flush() error {
	// upload vertice + meta to opengl
	GlBindBuffer(ARRAY_BUFFER, b.glBuffer)
	GlBufferData(b.glType, allVerts, STATIC_DRAW)
	// FIXME: rebind shader attrs
	b.buf = b.buf[:0]
	return nil
}
