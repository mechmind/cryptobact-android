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
    ID_FOOD
    TOTAL_IDS
)

var glSizeMap = map[C.GLuint]C.GLuint{
    C.GL_TRIANGLES: 3,
}

var mainSet = []objectSet{
    ID_BACTERIA: objectSet{
        glType: C.GL_LINE_LOOP,
    },
    ID_FOOD: objectSet{},
}

type objectSet struct {
    glBufferId C.GLuint
    glType C.GLenum
    glColor int
    objPattern []C.GLfloat
    vxs []C.GLfloat
    vxsBB []C.GLfloat
}

type Render struct {
    sets []objectSet
    posAttr C.GLuint
}

func newRender(posAttr C.GLuint) *Render {
    r := &Render{mainSet, posAttr}
    for _, set := range r.sets {
        set.glBufferId = GenBuffer()
    }
    return r
}

func (r *Render) UpdateSet(tag int, cx, cy, scale float32) {
    r.sets[tag].vxsBB = append(r.sets[tag].vxsBB,
        renderObject(r.sets[tag].objPattern, cx, cy, scale)...)
}

func (r *Render) SwapBB() {
    for _, set := range r.sets {
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
        C.glBindBuffer(C.GL_ARRAY_BUFFER, set.glBufferId)
        C.glVertexAttribPointer(r.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
        C.glDrawArrays(set.glType, 0, (C.GLsizei)(len(set.vxs) / 2))
    }
}

func (r *Render) Flush() {
    for _, set := range r.sets {
        set.vxsBB = set.vxsBB[:0]
    }
}


func renderObject(pattern []C.GLfloat, cx, cy, scale float32) []C.GLfloat {
    vexs := make([]C.GLfloat, len(pattern))
    for idx := 0; idx < len(pattern); idx += 2 {
        vexs[idx] = C.GLfloat(cx) + pattern[idx] * C.GLfloat(scale)
        vexs[idx + 1] = C.GLfloat(cy) + pattern[idx + 1] * C.GLfloat(scale)
    }
    return vexs
}
