package gl

// just wrappers around opengl C api

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"

import (
	"unsafe"
)

// opengl constants
const (
	STATIC_DRAW      = C.GL_STATIC_DRAW
	ARRAY_BUFFER     = C.GL_ARRAY_BUFFER
	CULL_FACE        = C.GL_CULL_FACE
	DEPTH_TEST       = C.GL_DEPTH_TEST
	COLOR_BUFFER_BIT = C.GL_COLOR_BUFFER_BIT
	DEPTH_BUFFER_BIT = C.GL_DEPTH_BUFFER_BIT
	TRUE             = C.GL_TRUE
	FALSE            = C.GL_FALSE
	INFO_LOG_LENGTH  = C.GL_INFO_LOG_LENGTH
	COMPILE_STATUS = C.GL_COMPILE_STATUS

	TRIANGLES = C.GL_TRIANGLES
	POINTS    = C.GL_POINTS

	FLOAT = C.GL_FLOAT

	VERTEX_SHADER   = C.GL_VERTEX_SHADER
	FRAGMENT_SHADER = C.GL_FRAGMENT_SHADER
)

func GlBindBuffer(buf uint, mode uint) error {
	C.glBindBuffer(C.GLenum(buf), C.GLuint(mode))
	return nil
}

func GlBufferData(glType uint, verts []Vertex, glMode uint) error {
	C.glBufferData(C.GLenum(glType),
		C.GLsizeiptr(len(verts)*int(unsafe.Sizeof(verts[0]))),
		unsafe.Pointer(&verts[0]), C.GLenum(glMode))
	return nil
}

func GlGenBuffer() (uint, error) {
	var buf C.GLuint
	C.glGenBuffers(1, &buf)
	return uint(buf), nil
}

func GlClearColor(r, g, b, a float32) error {
	C.glClearColor(C.GLclampf(r), C.GLclampf(g), C.GLclampf(b), C.GLclampf(a))
	return nil
}

func GlEnable(capa uint) error {
	C.glEnable(C.GLenum(capa))
	return nil
}

func GlUseProgram()              {}
func GlEnableVertexAttribArray() {}
func GlUniformMatrix4fv()        {}
func GlGetProgramiv()            {}
func GlGetProgramInfoLog()       {}

func GlShaderSource(handle uint, sources [][]string) error {
	sourcesC := make([]*C.GLchar, len(sources))
	sourceLengths := make([]C.GLint, len(sources))

	for idx, source := range sources {
		sourcesC[idx] = (*C.GLchar)(unsafe.Pointer(&source[0]))
		sourceLengths[idx] = C.GLint(len(source))
	}
	C.glShaderSource(C.GLenum(handle), C.GLint(len(sources)), &sourcesC[0], &sourceLengths[0])
	return nil
}

func GlCreateShader(shType uint) (uint, error) {
	handle := C.glCreateShader(C.GLenum(shType))
	return uint(handle), nil
}

func GlCompileShader(handle uint) error {
	C.glCompileShader(C.GLuint(handle))
	return nil
}

func GlGetShaderiv(handle uint, flag uint) (int, error) {
	var dest int
	C.glGetShaderiv(C.GLuint(handle), C.GLenum(flag), (*C.GLint)(&dest))
	return dest, nil
}

func GlGetShaderInfoLog()  {}
func GlGetProgramiv()      {}
func GlGetProgramInfoLog() {}
