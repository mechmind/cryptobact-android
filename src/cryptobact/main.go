package main

import (
	"cryptobact/engine"
	"cryptobact/gl"
	"cryptobact/ui"
	"log"
	"runtime"
)

const (
	X_COUNT        = 16
	Y_COUNT        = 24
	CROSS_HALFSIZE = 2.5

	STEP = 25.0
)

type game struct {
	width, height int
	mvp           []float32

	fieldScreen   *ui.FieldScreen
	presetScreen  *ui.PresetScreen
	currentScreen ui.UInteractive

	updater *Updater
}

var g game

var _ = engine.CALIBRATE_MS

func main() {
	runtime.GOMAXPROCS(2)
	g.updater = newUpdater()
	go g.updater.fetchUpdates()
	//go engine.Loop(g.updater)
}

//func updateCurrentBuffer(verts []C.GLfloat) {
//	C.glBufferData(C.GL_ARRAY_BUFFER,
//		C.GLsizeiptr(len(verts)*int(unsafe.Sizeof(verts[0]))),
//		unsafe.Pointer(&verts[0]), C.GL_STATIC_DRAW)
//}

func (game *game) resize(width, height int) {
	log.Println("now resize to ", width, height)
	game.width = width
	game.height = height

	//game.offsetX = float32(width-X_COUNT*STEP) / 2.0
	//game.offsetY = float32(height-Y_COUNT*STEP) / 2.0

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
	//	game.prog, _ = gl.CreateProgram(vertShaderSrcDef, fragShaderSrcDef)
	//	posAttrib, _ := gl.GlGetAttribLocation(game.prog, "vPosition")
	//	game.posAttr = posAttrib
	//	game.offsetUni, _ = gl.GlGetUniformLocation(game.prog, "offset")
	//	game.mvpUni, _ = gl.GlGetUniformLocation(game.prog, "mvp")
	//	game.colorUni, _ = gl.GlGetUniformLocation(game.prog, "color")
	//	gl.GlUseProgram(game.prog)
	//	gl.GlEnableVertexAttribArray(posAttrib)
	//	// transformation matrix
	//	gl.GlUniformMatrix4fv(game.mvpUni, 1, false, game.mvp)
	//
	//	game.gridBufId, _ = gl.GlGenBuffer()
	//	game.sliderLinesBufId, _ = gl.GlGenBuffer()
	//	game.sliderTriagsBufId, _ = gl.GlGenBuffer()
	//	// set up grid buffer
	//	game.verts = makeGridPoints(X_COUNT*STEP, Y_COUNT*STEP, STEP)
	//	gl.GlBindBuffer(gl.ARRAY_BUFFER, game.gridBufId)
	//	updateCurrentBuffer(game.verts)
	//

	//	// start engine
	//	game.render = newRender(posAttrib)
	//	game.updater.AttachRender(game.render)

	game.currentScreen = ui.NewHookerScreen(game)
	log.Println("screen: now hooker")
}

func (game *game) drawFrame() {
	gl.GlClear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	game.currentScreen.HandleDraw()
}

func (game *game) onTouch(action int, x, y float32) {
	game.currentScreen.HandleTouch(action, x, y)
}

func (g *game) Init(w, h int) {
	var fieldx float32 = STEP * X_COUNT
	var fieldy float32 = STEP * Y_COUNT

	fw := float32(w)
	fh := float32(h)
	if fieldx > fw || fieldy > fh {
		log.Println("screen: khooyovoyo razreshenie", fieldx, fw, fieldy, fh)
		return
	}

	hpad := (fw - fieldx) / 2
	bottompad := fh - fieldy - hpad

	bottomRect := ui.SimpleRect{0, 0, fw, bottompad}

	g.fieldScreen = ui.NewFieldScreen(g, bottomRect)
	g.fieldScreen.Init(w, h)
	//g.presetScreen = newPresetScreen(float32(fw * 2), 0, bottomRect)
	g.presetScreen = ui.NewPresetScreen(g, bottomRect)
	g.presetScreen.Init(w, h)

	g.updater.AttachField(g.fieldScreen.F)
	g.Switch(ui.ID_FIELD_SCREEN)
	log.Println("screen: now game")
}

func (g *game) Switch(id int) {
	switch id {
	case ui.ID_FIELD_SCREEN:
		g.currentScreen = g.fieldScreen
	case ui.ID_PRESET_SCREEN:
		g.currentScreen = g.presetScreen
	}
}

func (g *game) GetProjectionMatrix() []float32 {
	return g.mvp
}

func (g *game) GetSize() (h, w int, step float32) {
	return g.height, g.width, STEP
}
