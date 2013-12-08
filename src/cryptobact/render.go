package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
*/
import "C"
import "unsafe"

const (
	ID_BACTERIA = iota
	ID_EGG
	ID_FOOD
	TOTAL_IDS
)

var glSizeMap = map[C.GLuint]C.GLuint{
	C.GL_TRIANGLES: 3,
}

var mainSet = []*objectSet{
	ID_BACTERIA: &objectSet{
		glType:  C.GL_TRIANGLES,
		glColor: [3]C.GLfloat{1.0, 0.8, 0.4},
		objPattern: []C.GLfloat{
			0, -30, 0, 0, -8, -6,
			0, -30, 8, -6, 0, 0,
			8, -6, 12, -2, 0, 0,
			12, -2, 0, 6, 0, 0,
			0, 6, -12, -2, 0, 0,
			-12, -2, -8, -6, 0, 0,
		},
	},
	ID_EGG: &objectSet{
		glType:  C.GL_TRIANGLES,
		glColor: [3]C.GLfloat{0.1, 0.2, 0.8},
		objPattern: []C.GLfloat{
			0, 0, -5, -11, 5, -11,
			0, 0, 5, -11, 8, -7,
			0, 0, 8, -7, 8, 7,
			0, 0, 8, 7, 5, 11,
			0, 0, 5, 11, -5, 11,
			0, 0, -5, 11, -8, 7,
			0, 0, -8, 7, -5, -11,
		},
	},
	ID_FOOD: &objectSet{
		glType:  C.GL_TRIANGLES,
		glColor: [3]C.GLfloat{0.8, 0, 0.8},
		objPattern: []C.GLfloat{
			0, 0, 3, 5, -3, 5,
		},
	},
}

var colorSet = [][3]C.GLfloat{
	{1.0, 0, 0},
	{0, 1.0, 0},
	{0.9, 0.9, 0.9},
	{0, 0, 1.0},
}

var defaultColor = [3]C.GLfloat{1.0, 0.8, 0}

type splat struct {
	length  int
	glColor [3]C.GLfloat
}

type objectSet struct {
	glBufferId C.GLuint
	glType     C.GLenum
	glColor    [3]C.GLfloat
	objPattern []C.GLfloat
	vxs        []C.GLfloat
	vxsBB      []C.GLfloat
	splats     []splat
}

type Render struct {
	sets    []*objectSet
	posAttr C.GLuint
}

func newRender(posAttr C.GLuint) *Render {
	r := &Render{mainSet, posAttr}
	for _, set := range r.sets {
		set.glBufferId = GenBuffer()
		set.vxs = []C.GLfloat{}
		set.vxsBB = []C.GLfloat{}
	}
	return r
}

func (r *Render) UpdateSet(tag int, cx, cy, scale float32) int {
	vxs := renderObject(r.sets[tag].objPattern, cx, cy, scale)
	r.sets[tag].vxsBB = append(r.sets[tag].vxsBB, vxs...)
	return len(vxs)
	//log.Println("handled update set", tag, "now have", len(r.sets[tag].vxsBB), "vecs")
}

func (r *Render) ClearSplat(tag int, id int) {
	set := r.sets[tag]
	if len(set.splats) > id {
		set.splats[id].length = 0
	}
}

func (r *Render) UpdateSplat(tag int, id int, count int, color [3]C.GLfloat) {
	set := r.sets[tag]
	if len(set.splats) <= id {
		splats := make([]splat, id+1)
		copy(splats, set.splats)
		set.splats = splats
	}
	set.splats[id].length += count
	set.splats[id].glColor = color
}

func (r *Render) SwapBB() {
	for _, set := range r.sets {
		if len(set.vxsBB) == 0 {
			// no points
			continue
		}
		set.vxs, set.vxsBB = set.vxsBB, set.vxs
		set.vxsBB = set.vxsBB[:0]
		// set gl buffer
		C.glBindBuffer(C.GL_ARRAY_BUFFER, set.glBufferId)
		C.glVertexAttribPointer(r.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
		updateCurrentBuffer(set.vxs)
	}
}

func (r *Render) RenderAll() {
	for _, set := range r.sets {
		//log.Println("set", id, "has", len(set.vxs), "points")
		if len(set.vxs) == 0 {
			continue
		}

		if set.splats != nil {
			//log.Println("rendering splats", set.splats, "on vxs", len(set.vxs))
			var current int
			C.glBindBuffer(C.GL_ARRAY_BUFFER, set.glBufferId)
			for _, sp := range set.splats {
				if sp.length == 0 {
					continue
				}
				C.glVertexAttribPointer(r.posAttr, 2, C.GL_FLOAT, C.GL_FALSE,
					0, unsafe.Pointer(uintptr(0)))
				vxset := set.vxs[current : current+sp.length]
				updateCurrentBuffer(vxset)
				C.glVertexAttribPointer(r.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
				C.glUniform3f(C.GLint(g.colorUni), sp.glColor[0], sp.glColor[1], sp.glColor[2])
				C.glDrawArrays(set.glType, 0, (C.GLsizei)(len(vxset)/2))
				current += sp.length
				//set.splats[id].length = 0
			}
		} else {
			C.glBindBuffer(C.GL_ARRAY_BUFFER, set.glBufferId)
			C.glVertexAttribPointer(r.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
			C.glUniform3f(C.GLint(g.colorUni), set.glColor[0], set.glColor[1], set.glColor[2])
			C.glDrawArrays(set.glType, 0, (C.GLsizei)(len(set.vxs)/2))
		}
	}
}

func (r *Render) Flush() {
	for _, set := range r.sets {
		set.vxsBB = set.vxsBB[:0]
	}
}

func renderObject(pattern []C.GLfloat, cx, cy, scale float32) []C.GLfloat {
	vexs := make([]C.GLfloat, len(pattern))
	//log.Println("render: coords are", cx, cy)
	for idx := 0; idx < len(pattern); idx += 2 {
		vexs[idx] = C.GLfloat(cx)*STEP + pattern[idx]*C.GLfloat(scale)
		vexs[idx+1] = C.GLfloat(cy)*STEP + pattern[idx+1]*C.GLfloat(scale)
	}
	//log.Println("render: resulting coord set is", vexs)
	return vexs
}
