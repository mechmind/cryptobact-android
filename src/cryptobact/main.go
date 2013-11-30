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
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"unsafe"
)

var time float64

type game struct {
	prog                C.GLuint
	width, height       int
	offsetUni, colorUni int

	mu               sync.Mutex // Protects offsetX, offsetY
	offsetX, offsetY float32

	touching         bool
	touchX, touchY   float32
}

var g game

const vertShaderSrcDef = `
	uniform vec2 offset;
	attribute vec4 vPosition;

	void main() {
		gl_Position = vec4(vPosition.xy+offset, vPosition.zw);
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
    go runServer()
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

func (game *game) resize(width, height int) {
	game.width = width
	game.height = height
	C.glViewport(0, 0, C.GLsizei(width), C.GLsizei(height))
}

func (game *game) initGL() {
	log.Printf("GL_VERSION: %v GL_RENDERER: %v GL_VENDOR %v\n",
		GetString(C.GL_VERSION), GetString(C.GL_RENDERER), GetString(C.GL_VENDOR))
	log.Printf("GL_EXTENSIONS: %v\n", GetString(C.GL_EXTENSIONS))
	C.glClearColor(0.0, 0.0, 0.0, 1.0)
	C.glEnable(C.GL_CULL_FACE)
	C.glEnable(C.GL_DEPTH_TEST)

	game.prog = createProgram(vertShaderSrcDef, fragShaderSrcDef)
	posAttrib := attribLocation(game.prog, "vPosition")
	game.offsetUni = uniformLocation(game.prog, "offset")
	game.colorUni = uniformLocation(game.prog, "color")
	C.glUseProgram(game.prog)
	C.glEnableVertexAttribArray(C.GLuint(posAttrib))

	vertVBO := GenBuffer()
	checkGLError()
	C.glBindBuffer(C.GL_ARRAY_BUFFER, vertVBO)
	verts := []float32{.0, 0.5, -0.5, -0.5, 0.5, -0.5}
	C.glBufferData(C.GL_ARRAY_BUFFER, C.GLsizeiptr(len(verts)*int(unsafe.Sizeof(verts[0]))), unsafe.Pointer(&verts[0]), C.GL_STATIC_DRAW)
	C.glVertexAttribPointer(C.GLuint(posAttrib), 2, C.GL_FLOAT, C.GL_FALSE, 0, unsafe.Pointer(uintptr(0)))
}

func (game *game) drawFrame() {
	time += .05
	color := (C.GLclampf(math.Sin(time)) + 1) * .5

	game.mu.Lock()
	offX := game.offsetX
	offY := game.offsetY
	game.mu.Unlock()
	C.glUniform2f(C.GLint(game.offsetUni), C.GLfloat(offX), C.GLfloat(offY))
	C.glUniform3f(C.GLint(game.colorUni), 1.0, C.GLfloat(color), 0)
	C.glClear(C.GL_COLOR_BUFFER_BIT | C.GL_DEPTH_BUFFER_BIT)

	C.glUseProgram(game.prog)
	C.glDrawArrays(C.GL_TRIANGLES, 0, 3)
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
		game.mu.Lock()
		game.offsetX += 2 * (x - game.touchX) / float32(game.width)
		game.offsetY += 2 * -(y - game.touchY) / float32(game.height)
		game.mu.Unlock()
		game.touchX, game.touchY = x, y
	}
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
