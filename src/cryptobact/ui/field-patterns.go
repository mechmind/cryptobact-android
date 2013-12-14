package ui

import (
	"cryptobacteria/gl"
)

const (
	ID_BACTERIA = iota
	ID_EGG
	ID_FOOD
	TOTAL_IDS
)

type template struct {
	color    gl.Color
	vertices []gl.Vertex
}

var mainSet = []*template {
	ID_BACTERIA: &template{
		color: gl.Color{1.0, 0.8, 0.4, 1.0},
		vertices: []gl.Vertex{
			{0, -30}, {0, 0}, {-8, 6},
		},
	},
}

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
