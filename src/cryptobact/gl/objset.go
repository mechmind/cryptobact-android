package gl

type shaderBinder func(set *ObjectSet) error

type ObjectSet struct {
	glColor    Color
	objPattern []Vertex
	vxs        []Vertex
	vxsBB      []Vertex
	//splats     []splat
}

func NewObjectSet(color Color, pattern []Vertex) *ObjectSet {
	return &ObjectSet{color, pattern, nil, nil}
}

func (s *ObjectSet) MakeData() []ColoredVertex {
	return []ColoredVertex{}
}

type Buffer struct {
	sets     []*ObjectSet
	glType   uint
	glBuffer uint
	binder   shaderBinder
}

func NewBuffer(glType uint, binder shaderBinder) *Buffer {
	glBuf, _ := GlGenBuffer() // FIXME: handle error
	return &Buffer{nil, glType, glBuf, binder}
}

// upload data to opengl
// TODO: implement incremental update
func (b *Buffer) UploadData() error {
	// make array of all vertices
	vertSets := make([][]ColoredVertex, 0)
	var count int
	for _, set := range b.sets {
		verts := set.MakeData()
		vertSets = append(vertSets, verts)
		count += len(verts)
	}

	allVerts := make([]ColoredVertex, count)
	count = 0
	for _, verts := range vertSets {
		count += copy(allVerts[count:], verts)
	}

	// upload it to opengl
	GlBindBuffer(ARRAY_BUFFER, b.glBuffer)
	GlBufferData(b.glType, allVerts, STATIC_DRAW)
	// FIXME: rebind shader attrs
	return nil
}
