package main


/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"

type Square struct {
    pos int
    size float32
}

func (sq Square) Coords() []C.GLfloat {
    coords := make([]C.GLfloat, 6)

    cx, cy := sq.Center()
    size := C.GLfloat(sq.size)
    coords[0] = cx - size
    coords[1] = cy - size
    coords[2] = cx + size
    coords[3] = cy - size
    coords[4] = cx + size
    coords[5] = cy + size
    return coords
}

func (sq Square) Center() (x, y C.GLfloat) {
    x = STEP / 2
    y = STEP / 2

    x += C.GLfloat(sq.pos % X_COUNT) * STEP
    y += C.GLfloat(sq.pos / X_COUNT) * STEP

    return
}
