package ui

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"

const (
	SLIDER_X_PAD    = 10
	SLIDER_X_BUTTON = 50
	SLIDER_X_MARGIN = 14
	SLIDER_X_CURSOR = 25
	SLIDER_Y_SIZE   = 50
	SLIDER_Y_PAD    = 126
	SLIDER_Y_MARGIN = 56
)

type Slider struct {
	value float32
}

func (s *Slider) Vertexes(cx, cy, sizex, sizey C.GLfloat) (lines, triangles []C.GLfloat) {
	// total:
	//   1 line
	//   3 triangles

	left_corner := cx + SLIDER_X_PAD
	right_corner := cx + sizex - SLIDER_X_PAD
	cursor_pos := (right_corner-left_corner)*C.GLfloat(s.value) + left_corner
	cursor_top := cy + SLIDER_Y_SIZE - 20
	cursor_bottom := cy + 20
	triangles = []C.GLfloat{
		// left (decrease) arrow
		left_corner, cy + SLIDER_Y_SIZE/2, left_corner + SLIDER_X_BUTTON, cy,
		left_corner + SLIDER_X_BUTTON, cy + SLIDER_Y_SIZE,
		// right (increase) arrow
		right_corner, cy + SLIDER_Y_SIZE/2, right_corner - SLIDER_X_BUTTON, cy + SLIDER_Y_SIZE,
		right_corner - SLIDER_X_BUTTON, cy,
		// top (cursor) arrow
		cursor_pos - SLIDER_X_CURSOR, cursor_top, cursor_pos, cursor_bottom,
		cursor_pos, cursor_top,
	}

	lines = []C.GLfloat{left_corner + SLIDER_X_MARGIN, cy + SLIDER_Y_SIZE/2,
		right_corner - SLIDER_X_MARGIN, cy + SLIDER_Y_SIZE/2,
	}
	return lines, triangles
}

func (s *Slider) Hitboxes() []SimpleRect {
	return nil
}
