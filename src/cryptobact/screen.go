package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"
import "cryptobact/gl"
import "log"
import "unsafe"

type UInteractive interface {
	Offset() (x, y float32)
	HandleTouch(action int, x, y float32)
	HandleDraw()
	HandleResize(w, h int)
}

type Screen struct {
	offsetX, offsetY float32
}

func (s *Screen) Offset() (x, y float32) {
	return s.offsetX, s.offsetY
}

type touchTracker struct {
	x, y     float32
	touching bool
}

func (t *touchTracker) Update(action int, x, y float32, ofx, ofy float32) (x1, y1 float32) {
	x = x + float32(ofx)
	y = float32(g.height) - y + float32(ofy)
	switch action {
	case C.AMOTION_EVENT_ACTION_UP:
		t.touching = false
		t.x, t.y = x, y
		return x, y
	case C.AMOTION_EVENT_ACTION_DOWN:
		t.touching = true
		t.x, t.y = x, y
		return x, y
	case C.AMOTION_EVENT_ACTION_MOVE:
		if !t.touching {
			break
		}
		t.x, t.y = x, y
		return x, y
	}
	return x, y
}

type simpleRect struct {
	x1, y1, x2, y2 float32
}

func (s simpleRect) In(x, y float32) bool {
	return (s.x1 <= x && x <= s.x2 && s.y1 <= y && y <= s.y2)
}

type gameScreen struct {
	Screen
	t          touchTracker
	bottomRect simpleRect
}

func newGameScreen(offx, offy float32, bottomRect simpleRect) *gameScreen {
	return &gameScreen{Screen{offx, offy}, touchTracker{}, bottomRect}
}

func (gs *gameScreen) HandleTouch(action int, x, y float32) {
	// open control screen, if clicked in bottom part

	rx, ry := gs.t.Update(action, x, y, gs.Screen.offsetX, gs.Screen.offsetY)
	log.Println("recv coord", x, y, "converted coords:", rx, ry)

	if action == C.AMOTION_EVENT_ACTION_UP && gs.bottomRect.In(rx, ry) {
		log.Println("screen: game screen throws to preset")
		g.activateScreen(g.presetScreen)
	} else if action == C.AMOTION_EVENT_ACTION_UP {
		log.Println("screen: EVENT UP", rx, ry)
	}
}

func (gs *gameScreen) HandleDraw() {
	// render grid
	gl.GlBindBuffer(gl.ARRAY_BUFFER, g.gridBufId)
	gl.GlVertexAttribPointer(g.posAttr, 2, gl.FLOAT, false, 0, 0)
	C.glDrawArrays(C.GL_POINTS, 0, (C.GLsizei)(len(g.verts)))
	// world
	if status := g.updater.isWorldUpdated(); status != nil {
		// apply bb to render
		//log.Println("applying new map")
		g.render.SwapBB()
		status <- struct{}{}
	}
	//log.Println("rendering map")
	g.render.RenderAll()
}

func (gs *gameScreen) HandleResize(w, h int) {}

type presetScreen struct {
	Screen
	t          touchTracker
	bottomRect simpleRect
}

func newPresetScreen(offx, offy float32, bottomRect simpleRect) *presetScreen {
	return &presetScreen{Screen{offx, offy}, touchTracker{}, bottomRect}
}

func (ps *presetScreen) HandleTouch(action int, x, y float32) {
	// live drag
	//    if ps.t.touching {
	//        ps.Screen.offsetX += float32(x - ps.t.x)
	//        ps.Screen.offsetY -= float32(y - ps.t.y)
	//    }

	rx, ry := ps.t.Update(action, x, y, ps.Screen.offsetX, ps.Screen.offsetY)
	log.Println("recv coords", x, y, "converted coords:", rx, ry)
	if action == C.AMOTION_EVENT_ACTION_UP && ps.bottomRect.In(rx, ry) {
		log.Println("screen: pres screen throws to game")
		g.activateScreen(g.gameScreen)
	} else if action == C.AMOTION_EVENT_ACTION_UP {
		log.Println("screen: pres EVENT UP", x, y)
	}
}

func (ps *presetScreen) HandleDraw() {
	g.sliderLineBuffer = g.sliderLineBuffer[:0]
	g.sliderTrBuffer = g.sliderTrBuffer[:0]
	sizex, sizey := g.width, g.height
	basex, basey := ps.Offset()
	for idx, slider := range g.sliders {
		ypos := basey + float32(sizey) - float32(idx)*(SLIDER_Y_SIZE+SLIDER_Y_MARGIN) - SLIDER_Y_PAD
		lines, triags := slider.Vertexes(basex, ypos, float32(sizex), float32(sizey))
		g.sliderLineBuffer = append(g.sliderLineBuffer, lines...)
		g.sliderTrBuffer = append(g.sliderTrBuffer, triags...)
	}

	// render sliders
	// lines
	C.glBindBuffer(C.GL_ARRAY_BUFFER, g.sliderLinesBufId)
	updateCurrentBuffer(g.sliderLineBuffer)
	C.glVertexAttribPointer(g.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
	C.glDrawArrays(C.GL_LINES, 0, (C.GLsizei)(len(g.sliderLineBuffer)/2))

	//log.Println("pres screen: writing lines", g.sliderLineBuffer)

	// triags
	C.glBindBuffer(C.GL_ARRAY_BUFFER, g.sliderTriagsBufId)
	updateCurrentBuffer(g.sliderTrBuffer)
	C.glVertexAttribPointer(g.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
	C.glDrawArrays(C.GL_TRIANGLES, 0, (C.GLsizei)(len(g.sliderTrBuffer)/2))

	//log.Println("pres screen: writing triags", g.sliderTrBuffer)
}

func (ps *presetScreen) HandleResize(w, h int) {}

type hookerScreen struct{}

func newHookerScreen() *hookerScreen {
	return &hookerScreen{}
}

func (hs *hookerScreen) Offset() (x, y float32) {
	return 0, 0
}
func (hs *hookerScreen) HandleTouch(action int, x, y float32) {}
func (hs *hookerScreen) HandleDraw()                          {}

func (hs *hookerScreen) HandleResize(w, h int) {
	g.createScreens(w, h)
}

func (g *game) activateScreen(s UInteractive) {
	g.currentScreen = s
	dx, dy := s.Offset()
	g.offsetX, g.offsetY = float32(dx), float32(dy)
}

func (g *game) createScreens(w, h int) {
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

	bottomRect := simpleRect{0, 0, fw, bottompad}

	g.gameScreen = newGameScreen(float32(hpad), float32(bottompad), bottomRect)
	//g.presetScreen = newPresetScreen(float32(fw * 2), 0, bottomRect)
	g.presetScreen = newPresetScreen(0, 0, bottomRect)

	g.activateScreen(g.gameScreen)
	log.Println("screen: now game")
}
