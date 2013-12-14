package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"

import (
	"cryptobact/engine"
	"cryptobact/gl"
	"log"
	"math"
	"runtime"
	"sync"
	"unsafe"
)

const (
	X_COUNT        = 16
	Y_COUNT        = 24
	CROSS_HALFSIZE = 2.5

	STEP = 25.0
)

var ticks float64

type game struct {
	prog                uint
	width, height       int
	offsetUni, colorUni int
	posAttr             int
	mvpUni              int
	mvp                 []float32

	mu               sync.Mutex // Protects offsetX, offsetY
	offsetX, offsetY float32

	// buffer ids
	gridBufId                           uint
	sliderLinesBufId, sliderTriagsBufId uint

	verts []gl.ColoredVertex

	updater *Updater
	render  *Render

	sliders          []Slider
	sliderLineBuffer []gl.ColoredVertex
	sliderTrBuffer   []gl.ColoredVertex

	gameScreen   *gameScreen
	presetScreen *presetScreen

	currentScreen UInteractive
}

var g game

const vertShaderSrcDef = `
	uniform vec2 offset;
    uniform mat4 mvp;
	attribute vec4 vPosition;

	void main() {
		gl_Position = mvp * vec4(vPosition.xy+offset, vPosition.zw);
	}
`

const fragShaderSrcDef = `
	precision mediump float;

	uniform vec3 color;

	void main() {
		gl_FragColor = vec4(color.xyz, 1.0);
	}
`

func main() {
	runtime.GOMAXPROCS(2)
	g.updater = newUpdater()
	go g.updater.fetchUpdates()
	go engine.Loop(g.updater)

	g.sliders = []Slider{
		Slider{0.5},
		Slider{0.5},
		Slider{0.5},
	}
}

func updateCurrentBuffer(verts []C.GLfloat) {
	C.glBufferData(C.GL_ARRAY_BUFFER,
		C.GLsizeiptr(len(verts)*int(unsafe.Sizeof(verts[0]))),
		unsafe.Pointer(&verts[0]), C.GL_STATIC_DRAW)
}

func (game *game) resize(width, height int) {
	log.Println("now resize to ", width, height)
	game.width = width
	game.height = height

	game.offsetX = float32(width-X_COUNT*STEP) / 2.0
	game.offsetY = float32(height-Y_COUNT*STEP) / 2.0

	gl.MakeProjectionMatrix(0, float32(width)-1, 0, float32(height)-1, 1.0, -1.0, game.mvp)
	gl.GlViewport(0, 0, width, height)

	game.currentScreen.HandleResize(width, height)
}

func (game *game) initGL() {
	log.Println("initializing gl")
	gl.GlClearColor(0, 0, 0, 1)
	gl.GlEnable(gl.CULL_FACE)
	gl.GlEnable(gl.DEPTH_TEST)

	game.mvp = make([]float32, 16)
	game.prog, _ = gl.CreateProgram(vertShaderSrcDef, fragShaderSrcDef)
	posAttrib, _ := gl.GlGetAttribLocation(game.prog, "vPosition")
	game.posAttr = posAttrib
	game.offsetUni, _ = gl.GlGetUniformLocation(game.prog, "offset")
	game.mvpUni, _ = gl.GlGetUniformLocation(game.prog, "mvp")
	game.colorUni, _ = gl.GlGetUniformLocation(game.prog, "color")
	gl.GlUseProgram(game.prog)
	gl.GlEnableVertexAttribArray(posAttrib)
	// transformation matrix
	gl.GlUniformMatrix4fv(game.mvpUni, 1, false, game.mvp)

	game.gridBufId, _ = gl.GlGenBuffer()
	game.sliderLinesBufId, _ = gl.GlGenBuffer()
	game.sliderTriagsBufId, _ = gl.GlGenBuffer()
	// set up grid buffer
	game.verts = makeGridPoints(X_COUNT*STEP, Y_COUNT*STEP, STEP)
	gl.GlBindBuffer(gl.ARRAY_BUFFER, game.gridBufId)
	updateCurrentBuffer(game.verts)

	// start engine
	game.render = newRender(posAttrib)
	game.updater.AttachRender(game.render)

	game.currentScreen = newHookerScreen()
	log.Println("screen: now hooker")
}

func (game *game) drawFrame() {
	ticks += .05
	color := float32(math.Sin(ticks)+1) * .5

	game.mu.Lock()
	offX := game.offsetX
	offY := game.offsetY
	game.mu.Unlock()
	// basic stuff
	gl.GlUniform2f(game.offsetUni, offX, offY)
	gl.GlUniform3f(game.colorUni, 1.0, color, 0)
	gl.GlUniformMatrix4fv(game.mvpUni, 1, false, game.mvp)
	gl.GlClear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.GlUseProgram(game.prog)

	game.currentScreen.HandleDraw()
}

func (game *game) onTouch(action int, x, y float32) {
	game.currentScreen.HandleTouch(action, x, y)
}

func makeGridPoints(llimX, llimY, lstep float32) []gl.ColoredVertix {
	limX, limY, step := C.GLfloat(llimX), C.GLfloat(llimY), C.GLfloat(lstep)
	data := make([]C.GLfloat, 0, int(math.Ceil(float64(limX)*float64(limY)/(float64(step)*float64(step))+4)*4))

	var nextX, nextY C.GLfloat
	for nextX = 0.0; nextX < limX+0.1; nextX += step {
		for nextY = 0.0; nextY < limY+0.1; nextY += step {
			data = append(data, nextX, nextY)
			//            data = append(data, nextX - CROSS_HALFSIZE, nextY, nextX + CROSS_HALFSIZE, nextY)
			//            data = append(data, nextX, nextY - CROSS_HALFSIZE, nextX, nextY + CROSS_HALFSIZE)
		}
	}
	return data
}
