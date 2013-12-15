package ui

//import "cryptobact/gl"
import "log"

const (
	ID_FIELD_SCREEN = iota
	ID_PRESET_SCREEN
)

const (
	AMOTION_EVENT_ACTION_UP = iota
	AMOTION_EVENT_ACTION_DOWN
	AMOTION_EVENT_ACTION_MOVE
)

const (
	STEP    = 25.0 // FIXME: unconst all
	X_COUNT = 48
	Y_COUNT = 60
)

type UInteractive interface {
	Init(w, h int)
	HandleTouch(action int, x, y float32)
	HandleDraw()
	HandleResize(w, h int)
}

type Screener interface {
	Init(h, w int)
	Switch(screenId int)
	GetProjectionMatrix() []float32
	GetSize() (h, w int, step float32)
}

type Updater interface {
	IsWorldUpdated() chan struct{}
}

type Screen struct {
	screener Screener
}

type touchTracker struct {
	x, y             float32
	touchx, touchy   float32
	screenx, screeny float32
	bounds           SimpleRect
	touching         bool
}

func fmax(v1, v2 float32) float32 {
	if v1 > v2 {
		return v1
	}
	return v2
}

func fmin(v1, v2 float32) float32 {
	if v1 < v2 {
		return v1
	}
	return v2
}

func (t *touchTracker) Update(action int, x, y float32) (x1, y1 float32) {
	switch action {
	case AMOTION_EVENT_ACTION_UP:
		t.touching = false
	case AMOTION_EVENT_ACTION_DOWN:
		t.touching = true
		t.touchx = x
		t.touchy = y
	case AMOTION_EVENT_ACTION_MOVE:
		if !t.touching {
			break
		}
		dx := t.x + (x - t.touchx)
		dy := t.y - (y - t.touchy)
		t.x = fmax(fmin(dx, t.bounds.X2-t.screenx), t.bounds.X1)
		t.y = fmax(fmin(dy, t.bounds.Y2-t.screeny), t.bounds.Y1)
		t.touchx = x
		t.touchy = y
	}
	return t.x, t.y
}

type SimpleRect struct {
	X1, Y1, X2, Y2 float32
}

func (s SimpleRect) In(x, y float32) bool {
	return (s.X1 <= x && x <= s.X2 && s.Y1 <= y && y <= s.Y2)
}

type FieldScreen struct {
	Screen
	F          *Field
	t          touchTracker
	u          Updater
	bottomRect SimpleRect
}

func NewFieldScreen(sc Screener, u Updater, bottomRect SimpleRect) *FieldScreen {
	return &FieldScreen{Screen{sc}, NewField(), touchTracker{}, u, bottomRect}
}

func (fs *FieldScreen) Init(w, h int) {
	err := fs.F.Init(fs.screener.GetProjectionMatrix())
	if err != nil {
		panic(err)
	}
	scrx := float32(w)
	scry := float32(h)
	fs.t.screenx = scrx
	fs.t.screeny = scry
	fs.t.bounds = SimpleRect{-2 * scrx, -2 * scry, 2 * scrx, 2 * scry}

	log.Println("screen: initialized viewport", w, h)
}

func (fs *FieldScreen) HandleTouch(action int, x, y float32) {
	// open control screen, if clicked in bottom part

	ox, oy := fs.t.Update(action, x, y)
	if action == AMOTION_EVENT_ACTION_UP && fs.bottomRect.In(x, y) {
		//log.Println("screen: game screen throws to preset")
		//fs.screener.Switch(ID_PRESET_SCREEN)
	} else if action == AMOTION_EVENT_ACTION_UP {
		//log.Println("screen: EVENT UP", x, y)
	}
	fs.F.offx = ox
	fs.F.offy = oy
	//log.Println("screen: offset is", ox, oy)
}

func (fs *FieldScreen) HandleDraw() {
	// render grid

	mvp := fs.screener.GetProjectionMatrix()
	//	gl.GlBindBuffer(gl.ARRAY_BUFFER, g.gridBufId)
	//	gl.GlVertexAttribPointer(g.posAttr, 2, gl.FLOAT, false, 0, 0)
	//	C.glDrawArrays(C.GL_POINTS, 0, (C.GLsizei)(len(g.verts)))
	//	// world
	if status := fs.u.IsWorldUpdated(); status != nil {
		// apply bb to render
		//log.Println("applying new map")
		fs.F.FlushAll()
		status <- struct{}{}
	}
	fs.F.Draw(mvp)
	//	//log.Println("rendering map")
	//	g.render.RenderAll()
}

func (fs *FieldScreen) HandleResize(w, h int) {}

type PresetScreen struct {
	Screen
	t          touchTracker
	bottomRect SimpleRect
}

func NewPresetScreen(s Screener, bottomRect SimpleRect) *PresetScreen {
	return &PresetScreen{Screen{s}, touchTracker{}, bottomRect}
}

func (gs *PresetScreen) Init(w, h int) {}
func (ps *PresetScreen) HandleTouch(action int, x, y float32) {
	// live drag
	//    if ps.t.touching {
	//        ps.Screen.offsetX += float32(x - ps.t.x)
	//        ps.Screen.offsetY -= float32(y - ps.t.y)
	//    }

	if action == AMOTION_EVENT_ACTION_UP && ps.bottomRect.In(x, y) {
		log.Println("screen: pres screen throws to game")
		ps.screener.Switch(ID_FIELD_SCREEN)
	} else if action == AMOTION_EVENT_ACTION_UP {
		log.Println("screen: pres EVENT UP", x, y)
	}
}

func (ps *PresetScreen) HandleDraw() {
	//	g.sliderLineBuffer = g.sliderLineBuffer[:0]
	//	g.sliderTrBuffer = g.sliderTrBuffer[:0]
	//	sizex, sizey := g.width, g.height
	//	basex, basey := ps.Offset()
	//	for idx, slider := range g.sliders {
	//		ypos := basey + float32(sizey) - float32(idx)*(SLIDER_Y_SIZE+SLIDER_Y_MARGIN) - SLIDER_Y_PAD
	//		lines, triags := slider.Vertexes(basex, ypos, float32(sizex), float32(sizey))
	//		g.sliderLineBuffer = append(g.sliderLineBuffer, lines...)
	//		g.sliderTrBuffer = append(g.sliderTrBuffer, triags...)
	//	}
	//
	//	// render sliders
	//	// lines
	//	C.glBindBuffer(C.GL_ARRAY_BUFFER, g.sliderLinesBufId)
	//	updateCurrentBuffer(g.sliderLineBuffer)
	//	C.glVertexAttribPointer(g.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
	//	C.glDrawArrays(C.GL_LINES, 0, (C.GLsizei)(len(g.sliderLineBuffer)/2))
	//
	//	//log.Println("pres screen: writing lines", g.sliderLineBuffer)
	//
	//	// triags
	//	C.glBindBuffer(C.GL_ARRAY_BUFFER, g.sliderTriagsBufId)
	//	updateCurrentBuffer(g.sliderTrBuffer)
	//	C.glVertexAttribPointer(g.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
	//	C.glDrawArrays(C.GL_TRIANGLES, 0, (C.GLsizei)(len(g.sliderTrBuffer)/2))
	//
	//	//log.Println("pres screen: writing triags", g.sliderTrBuffer)
}

func (ps *PresetScreen) HandleResize(w, h int) {}

type hookerScreen struct {
	screener Screener
}

func NewHookerScreen(s Screener) *hookerScreen {
	return &hookerScreen{s}
}

func (hs *hookerScreen) Init(w, h int)                        {}
func (hs *hookerScreen) HandleTouch(action int, x, y float32) {}
func (hs *hookerScreen) HandleDraw()                          {}

func (hs *hookerScreen) HandleResize(w, h int) {
	hs.screener.Init(w, h)
}

//func (g *game) createScreens(w, h int) {
//	var fieldx float32 = STEP * X_COUNT
//	var fieldy float32 = STEP * Y_COUNT
//
//	fw := float32(w)
//	fh := float32(h)
//	if fieldx > fw || fieldy > fh {
//		log.Println("screen: khooyovoyo razreshenie", fieldx, fw, fieldy, fh)
//		return
//	}
//
//	hpad := (fw - fieldx) / 2
//	bottompad := fh - fieldy - hpad
//
//	bottomRect := SimpleRect{0, 0, fw, bottompad}
//
//	g.FieldScreen = newGameScreen(float32(hpad), float32(bottompad), bottomRect)
//	//g.PresetScreen = newPresetScreen(float32(fw * 2), 0, bottomRect)
//	g.PresetScreen = newPresetScreen(0, 0, bottomRect)
//
//	g.activateScreen(g.FieldScreen)
//	log.Println("screen: now game")
//}
