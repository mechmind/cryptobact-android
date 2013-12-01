package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2

void set_ortho_proj(GLfloat *matrix, GLfloat left, GLfloat right,
        GLfloat bottom, GLfloat top, GLfloat near, GLfloat far);

*/
import "C"

import (
	"cryptobact/engine"
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"unsafe"
)


const (
    X_COUNT = 16
    Y_COUNT = 24

    STEP = 25.0
)

var ticks float64

type game struct {
	prog                C.GLuint
	width, height       int
	offsetUni, colorUni int
    posAttr             C.GLuint
    mvpUni              int
    mvp                 []C.GLfloat

	mu               sync.Mutex // Protects offsetX, offsetY
	offsetX, offsetY float32

	touching         bool
	touchX, touchY   float32

    // buffer ids
    gridBufId C.GLuint

    verts []C.GLfloat

    updater *Updater
    render *Render
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
}

func GetShaderInfoLog(shader C.GLuint) string {
	var logLen C.GLint
	C.glGetShaderiv(shader, C.GL_INFO_LOG_LENGTH, &logLen)
	var c C.GLchar
	logLenBytes := int(logLen) * int(unsafe.Sizeof(c))
	log := C.malloc(C.size_t(logLenBytes))
	if log == nil {
		panic("Failed to allocate shader log buffer")
	}
	defer C.free(log)
	C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(logLen), (*C.GLsizei)(unsafe.Pointer(nil)), (*C.GLchar)(log))
	return string(C.GoBytes(log, C.int(logLenBytes)))
}

func GetProgramInfoLog(program C.GLuint) string {
	var logLen C.GLint
	C.glGetProgramiv(program, C.GL_INFO_LOG_LENGTH, &logLen)
	var c C.GLchar
	logLenBytes := int(logLen) * int(unsafe.Sizeof(c))
	log := C.malloc(C.size_t(logLenBytes))
	if log == nil {
		panic("Failed to allocate shader log buffer")
	}
	defer C.free(log)
	C.glGetProgramInfoLog(C.GLuint(program), C.GLsizei(logLen), (*C.GLsizei)(unsafe.Pointer(nil)), (*C.GLchar)(log))
	return string(C.GoBytes(log, C.int(logLenBytes)))
}

func loadShader(shaderType C.GLenum, source string) C.GLuint {
	handle := C.glCreateShader(shaderType)
	if handle == 0 {
		panic(fmt.Errorf("Failed to create shader of type %v", shaderType))
	}
	sourceC := C.CString(source)
	defer C.free(unsafe.Pointer(sourceC))
	C.glShaderSource(handle, 1, (**C.GLchar)(unsafe.Pointer(&sourceC)), (*C.GLint)(unsafe.Pointer(nil)))
	C.glCompileShader(handle)
	var compiled C.GLint
	C.glGetShaderiv(handle, C.GL_COMPILE_STATUS, &compiled)
	if compiled != C.GL_TRUE {
		log := GetShaderInfoLog(handle)
		panic(fmt.Errorf("Failed to compile shader: %v, shader: %v", log, source))
	}
	return handle
}

func GenBuffer() C.GLuint {
	var buf C.GLuint
	C.glGenBuffers(1, &buf)
	return C.GLuint(buf)
}

func checkGLError() {
	if glErr := C.glGetError(); glErr != C.GL_NO_ERROR {
		panic(fmt.Errorf("C.gl error: %v", glErr))
	}
}

func createProgram(vertShaderSrc string, fragShaderSrc string) C.GLuint {
	vertShader := loadShader(C.GL_VERTEX_SHADER, vertShaderSrc)
	fragShader := loadShader(C.GL_FRAGMENT_SHADER, fragShaderSrc)
	prog := C.glCreateProgram()
	if prog == 0 {
		panic("Failed to create shader program")
	}
	C.glAttachShader(prog, vertShader)
	checkGLError()
	C.glAttachShader(prog, fragShader)
	checkGLError()
	C.glLinkProgram(prog)
	var linkStatus C.GLint
	C.glGetProgramiv(prog, C.GL_LINK_STATUS, &linkStatus)
	if linkStatus != C.GL_TRUE {
		log := GetProgramInfoLog(prog)
		panic(fmt.Errorf("Failed to link program: %v", log))
	}
	return prog
}

func attribLocation(prog C.GLuint, name string) int {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	attrib := int(C.glGetAttribLocation(C.GLuint(prog), (*C.GLchar)(unsafe.Pointer(nameC))))
	checkGLError()
	if attrib == -1 {
		panic(fmt.Errorf("Failed to find attrib position for %v", name))
	}
	return attrib
}

func uniformLocation(prog C.GLuint, name string) int {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	attrib := int(C.glGetUniformLocation(C.GLuint(prog), (*C.GLchar)(unsafe.Pointer(nameC))))
	checkGLError()
	if attrib == -1 {
        panic(fmt.Errorf("Failed to find attrib position for %v", name))
	}
	return attrib
}

func GetString(name C.GLenum) string {
	val := C.glGetString(C.GLenum(name))
	return C.GoString((*C.char)(unsafe.Pointer(val)))
}

func updateCurrentBuffer(verts []C.GLfloat) {
	C.glBufferData(C.GL_ARRAY_BUFFER,
        C.GLsizeiptr(len(verts)*int(unsafe.Sizeof(verts[0]))),
        unsafe.Pointer(&verts[0]), C.GL_STATIC_DRAW)
}


func (game *game) resize(width, height int) {
	game.width = width
	game.height = height

    game.offsetX = float32(width - X_COUNT * STEP) / 2.0
    game.offsetY = float32(height - Y_COUNT * STEP) / 2.0

    C.set_ortho_proj((*C.GLfloat)(unsafe.Pointer(&game.mvp[0])), 0, C.GLfloat(width - 1),
        0, C.GLfloat(height - 1), 1.0, -1.0)
	C.glViewport(0, 0, C.GLsizei(width), C.GLsizei(height))
}

func (game *game) initGL() {
	log.Printf("GL_VERSION: %v GL_RENDERER: %v GL_VENDOR %v\n",
		GetString(C.GL_VERSION), GetString(C.GL_RENDERER), GetString(C.GL_VENDOR))
	log.Printf("GL_EXTENSIONS: %v\n", GetString(C.GL_EXTENSIONS))
	C.glClearColor(0.0, 0.0, 0.0, 1.0)
	C.glEnable(C.GL_CULL_FACE)
	C.glEnable(C.GL_DEPTH_TEST)

    game.mvp = make([]C.GLfloat, 16)
	game.prog = createProgram(vertShaderSrcDef, fragShaderSrcDef)
	posAttrib := attribLocation(game.prog, "vPosition")
    game.posAttr = C.GLuint(posAttrib)
	game.offsetUni = uniformLocation(game.prog, "offset")
    game.mvpUni = uniformLocation(game.prog, "mvp")
	game.colorUni = uniformLocation(game.prog, "color")
	C.glUseProgram(game.prog)
	C.glEnableVertexAttribArray(C.GLuint(posAttrib))
    // transformation matrix
    C.glUniformMatrix4fv(C.GLint(game.mvpUni), 1, C.GL_FALSE,
        (*C.GLfloat)(unsafe.Pointer(&game.mvp[0])))

	game.gridBufId = GenBuffer()
	checkGLError()
    // set up grid buffer
    game.verts = makeGridPoints(X_COUNT * STEP, Y_COUNT * STEP, STEP)
    C.glBindBuffer(C.GL_ARRAY_BUFFER, game.gridBufId)
    updateCurrentBuffer(game.verts)

    // start engine
    game.render = newRender(C.GLuint(posAttrib))
    game.updater = newUpdater(game.render)
    go game.updater.fetchUpdates()

    go engine.Loop(game.updater)
}

func (game *game) drawFrame() {
	ticks += .05
	color := (C.GLclampf(math.Sin(ticks)) + 1) * .5

	game.mu.Lock()
	offX := game.offsetX
	offY := game.offsetY
	game.mu.Unlock()
    // basic stuff
	C.glUniform2f(C.GLint(game.offsetUni), C.GLfloat(offX), C.GLfloat(offY))
	C.glUniform3f(C.GLint(game.colorUni), 1.0, C.GLfloat(color), 0)
    C.glUniformMatrix4fv(C.GLint(game.mvpUni), 1, C.GL_FALSE,
        (*C.GLfloat)(unsafe.Pointer(&game.mvp[0])))
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT)
	C.glUseProgram(game.prog)
    // grid
    C.glBindBuffer(C.GL_ARRAY_BUFFER, game.gridBufId)
	C.glVertexAttribPointer(game.posAttr, 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
    C.glDrawArrays(C.GL_LINES, 0, (C.GLsizei)(len(game.verts)))
    // world
    if status := game.updater.isWorldUpdated(); status != nil {
        // apply bb to render
        log.Println("applying new map")
        game.render.SwapBB()
        status <- struct{}{}
    }
    log.Println("rendering map")
    game.render.RenderAll()
}

func (game *game) onTouch(action int, x, y float32) {
	switch action {
	case C.AMOTION_EVENT_ACTION_UP:
		game.touching = false
	case C.AMOTION_EVENT_ACTION_DOWN:
		game.touching = true
		game.touchX, game.touchY = x, y
	case C.AMOTION_EVENT_ACTION_MOVE:
		if !game.touching {
			break
		}
        // ignore touch, BE UNTOUCHEABLE
//		game.mu.Lock()
//		game.offsetX += (x - game.touchX)
//		game.offsetY += -(y - game.touchY)
//		game.mu.Unlock()
		game.touchX, game.touchY = x, y
	}
}

func makeGridPoints(llimX, llimY, lstep float32) []C.GLfloat {
    limX, limY, step := C.GLfloat(llimX), C.GLfloat(llimY), C.GLfloat(lstep)
    data := make([]C.GLfloat, 0, int(math.Ceil(float64(limX) * float64(limY) / (float64(step) * float64(step)) + 4) * 4))

    var nextX, nextY C.GLfloat
    for nextX = 0.0 ; nextX < limX + 0.1 ; nextX += step {
        for nextY = 0.0; nextY < limY + 0.1; nextY += step {
            data = append(data, nextX, 0.0, nextX, limY)
            data = append(data, 0.0, nextY, limX, nextY)
        }
    }
    return data
}

// Use JNI_OnLoad to ensure that the go runtime is initialized at a predictable time,
// namely at System.loadLibrary()
//export JNI_OnLoad
func JNI_OnLoad(vm *C.JavaVM, reserved unsafe.Pointer) C.jint {
	return C.JNI_VERSION_1_6
}

//export Java_net_goandroid_cryptobact_Engine_drawFrame
func Java_net_goandroid_cryptobact_Engine_drawFrame(env *C.JNIEnv, clazz C.jclass) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: drawFrame: %v\n", err)
		}
	}()
	g.drawFrame()
}

//export Java_net_goandroid_cryptobact_Engine_init
func Java_net_goandroid_cryptobact_Engine_init(env *C.JNIEnv, clazz C.jclass) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: init: %v\n", err)
		}
	}()
	g.initGL()
}

//export Java_net_goandroid_cryptobact_Engine_resize
func Java_net_goandroid_cryptobact_Engine_resize(env *C.JNIEnv, clazz C.jclass, width, height C.jint) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: resize: %v\n", err)
		}
	}()
	g.resize(int(width), int(height))
}

//export Java_net_goandroid_cryptobact_Engine_onTouch
func Java_net_goandroid_cryptobact_Engine_onTouch(env *C.JNIEnv, clazz C.jclass, action C.jint, x, y C.jfloat) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("panic: resize: %v\n", err)
		}
	}()
	g.onTouch(int(action), float32(x), float32(y))
}
