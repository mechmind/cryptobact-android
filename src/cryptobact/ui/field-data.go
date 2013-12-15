package ui

import (
	"cryptobact/gl"
)

const (
	ID_BACTERIA_BODY = iota
	ID_BACTERIA_ORG
	ID_EGG
	ID_FOOD
	ID_MARKUP
	TOTAL_IDS
)

type template struct {
	verts  []float32
	glType uint
	color  [3]byte
}

var mainSet = []*template{
	ID_BACTERIA_BODY: &template{
		verts: []float32{
			0, -30, 0, 0, -8, 6,
		},
		glType: uint(gl.TRIANGLES),
		color:  [3]byte{},
	},
	ID_BACTERIA_ORG: &template{
		verts: []float32{
			0, -30, 0, 0, -8, 6,
		},
		glType: uint(gl.TRIANGLES),
		color:  [3]byte{},
	},
	ID_EGG: &template{
		verts: []float32{
			0, -30, 0, 0, -8, 6,
		},
		glType: uint(gl.TRIANGLES),
		color:  [3]byte{},
	},
	ID_FOOD: &template{
		verts: []float32{
			0, -30, 0, 0, -8, 6,
		},
		glType: uint(gl.TRIANGLES),
		color:  [3]byte{},
	},
	ID_MARKUP: &template{
		verts: []float32{
			1, 1,
		},
		glType: uint(gl.POINTS),
		color:  [3]byte{250, 160, 0},
	},
}

//var bacteriaVXShader = `
//    uniform vec2 offset;
//    uniform mat4 mvp;
//    attribute vec4 position;
//	attribute vec4 color;
//
//    void main() {
//        gl_Position = mvp * vec4(position.xy+offset, position.zw);
//    }
//`

var bacteriaVXShader = `
    uniform vec2 offset;
    uniform mat4 mvp;
    attribute vec4 position;

    void main() {
        gl_Position = vec4(position.xy+offset, position.zw);
    }
`
var bacteriaFragShader = `
    precision mediump float;

    uniform vec3 color;

    void main() {
        gl_FragColor = vec4(color.xyz, 1.0);
    }
`

/*
var mainSet = []*gl.ObjectSet{
	ID_BACTERIA: &objectSet{
		glType:  gl.TRIANGLES,
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
*/
