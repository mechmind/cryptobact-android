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
-4.8, -5.5, -4.0, -5.9, -5.5, -4.8,
-5.5, -4.8, -4.0, -5.9, -5.3, -4.0,
-5.3, -4.0, -4.0, -5.9, -3.6, -6.0,
-3.6, -6.0, -5.5, -3.2, -5.3, -4.0,
-5.5, -3.2, -3.6, -6.0, -3.0, -6.0,
-3.0, -6.0, -5.2, -2.5, -5.5, -3.2,
-5.2, -2.5, -3.0, -6.0, -2.2, -5.8,
-2.2, -5.8, -4.4, -1.6, -5.2, -2.5,
-4.4, -1.6, -2.2, -5.8, -1.6, -5.4,
-1.6, -5.4, -3.5, -1.0, -4.4, -1.6,
-3.5, -1.0, -1.6, -5.4, -1.0, -4.8,
-1.0, -4.8, -2.7, -0.5, -3.5, -1.0,
-2.7, -0.5, -1.0, -4.8, -0.6, -4.2,
-0.6, -4.2, -1.0, 1.0, -2.7, -0.5,
-0.6, -4.2, -0.4, -3.5, -1.0, 1.0,
-1.0, 1.0, -0.4, -3.5, -0.4, 1.3,
-0.4, 1.3, -0.4, -3.5, -0.2, -3.0,
-0.2, -3.0, 0.2, 2.2, -0.4, 1.3,
-0.2, -3.0, 0.8, -2.3, 0.2, 2.2,
0.2, 2.2, 0.8, -2.3, 0.5, 3.0,
0.5, 3.0, 0.8, -2.3, 2.1, -1.3,
2.1, -1.3, 0.9, 4.0, 0.5, 3.0,
2.1, -1.3, 4.5, -0.1, 0.9, 4.0,
0.9, 4.0, 4.5, -0.1, 2.0, 5.4,
2.0, 5.4, 4.5, -0.1, 5.0, 0.5,
5.0, 0.5, 2.5, 5.8, 2.0, 5.4,
5.0, 0.5, 5.8, 2.2, 2.5, 5.8,
2.5, 5.8, 5.8, 2.2, 3.0, 6.0,
3.0, 6.0, 5.8, 2.2, 5.8, 4.0,
5.8, 4.0, 3.5, 6.0, 3.0, 6.0,
3.5, 6.0, 5.8, 4.0, 5.1, 5.4,
5.1, 5.4, 4.5, 5.8, 3.5, 6.0,
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

var gridVXShader = `
    uniform vec2 offset;
    uniform mat4 mvp;
    attribute vec4 position;

    void main() {
        gl_Position = mvp * vec4(position.xy+offset, position.zw);
    }
`
var gridFragShader = `
    precision mediump float;

    uniform vec3 color;

    void main() {
        gl_FragColor = vec4(color.xyz, 1.0);
    }
`

var bactVXShader = `
    uniform vec2 offset;
    uniform mat4 mvp;
    attribute vec4 position;
	attribute vec4 lcolor;

	varying vec4 rcolor;

    void main() {
        gl_Position = mvp * vec4(position.xy+offset, position.zw);
		rcolor = lcolor.xyzw;
    }
`

var bactFragShader = `
    precision mediump float;

    varying vec4 rcolor;

	uniform vec3 color;

    void main() {
        gl_FragColor = vec4(rcolor.xyz, 1.0);
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
