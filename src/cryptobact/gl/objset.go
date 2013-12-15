package gl

import (
	"log"
)

type ShaderBinder func(set *Buffer) error

type Buffer struct {
	GlType   uint
	GlBuffer uint
	Buf      []float32
	BufLen   int
	Binder   ShaderBinder
}

func NewBuffer(glType uint, binder ShaderBinder) *Buffer {
	glBuf, _ := GlGenBuffer() // FIXME: handle error
	log.Println("gl: allocated buffer", glBuf)
	ErrPanic()
	return &Buffer{glType, glBuf, nil, 0, binder}
}

func (b *Buffer) Append(data []float32) {
	b.Buf = append(b.Buf, data...)
}

// flush data to opengl
// TODO: implement incremental update
func (b *Buffer) Flush() error {
	// upload vertice + meta to opengl
	GlBindBuffer(ARRAY_BUFFER, b.GlBuffer)
	ErrPanic()
	if len(b.Buf) == 0 {
		return nil
	}
	GlBufferData(ARRAY_BUFFER, b.Buf, STATIC_DRAW)
	ErrPanic()
	// FIXME: rebind shader attrs
	b.BufLen = len(b.Buf)
	b.Buf = b.Buf[:0]
	return nil
}
