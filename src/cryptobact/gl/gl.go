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
	"log"
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
	COMPILE_STATUS   = C.GL_COMPILE_STATUS
	LINK_STATUS      = C.GL_LINK_STATUS

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

func GlBufferData(glType uint, verts []float32, glMode uint) error {
	log.Printf("gl: loading %d vertices into %d", len(verts), glType)
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

func GlClear(mask C.GLbitfield) error {
	C.glClear(mask)
	return nil
}

func GlEnable(capa uint) error {
	C.glEnable(C.GLenum(capa))
	return nil
}

func GlCreateProgram() (uint, error) {
	handle := C.glCreateProgram()
	return uint(handle), nil
}

func GlLinkProgram(handle uint) error {
	C.glLinkProgram(C.GLuint(handle))
	return nil
}

func GlGetProgramiv(handle uint, flag uint) (int, error) {
	var dest int
	C.glGetProgramiv(C.GLuint(handle), C.GLenum(flag), (*C.GLint)(unsafe.Pointer(&dest)))
	return dest, nil
}

func GlGetProgramInfoLog(prog uint) (string, error) {
	logLen, _ := GlGetProgramiv(prog, INFO_LOG_LENGTH)
	if logLen == 0 {
		return "", nil
	}
	buf := make([]byte, logLen)
	C.glGetProgramInfoLog(C.GLuint(prog), C.GLsizei(logLen), (*C.GLsizei)(unsafe.Pointer(&logLen)),
		(*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf), nil
}

func GlUseProgram(handle uint) error {
	C.glUseProgram(C.GLuint(handle))
	return nil
}

func GlShaderSource(handle uint, sources []string) error {
	sourcesC := make([]*C.GLchar, len(sources))
	sourceLengths := make([]C.GLint, len(sources))

	for idx, source := range sources {
		byteSrc := []byte(source)
		sourcesC[idx] = (*C.GLchar)(unsafe.Pointer(&byteSrc[0]))
		sourceLengths[idx] = C.GLint(len(byteSrc))
	}
	log.Println("gl: loading shader into handle", handle, "count is", len(sources))
	C.glShaderSource(C.GLuint(handle), C.GLsizei(len(sources)),
		(**C.GLchar)(unsafe.Pointer(&sourcesC[0])), (*C.GLint)(unsafe.Pointer(&sourceLengths[0])))
	return nil
}

func GlCreateShader(shType uint) (uint, error) {
	handle := C.glCreateShader(C.GLenum(shType))
	log.Println("gl: created shader", handle)
	return uint(handle), nil
}

func GlCompileShader(handle uint) error {
	C.glCompileShader(C.GLuint(handle))
	return nil
}

func GlAttachShader(prog, shader uint) error {
	C.glAttachShader(C.GLuint(prog), C.GLuint(shader))
	return nil
}

func GlGetShaderiv(handle uint, flag uint) (int, error) {
	var dest int
	C.glGetShaderiv(C.GLuint(handle), C.GLenum(flag), (*C.GLint)(unsafe.Pointer(&dest)))
	return dest, nil
}

func GlGetShaderInfoLog(shader uint) (string, error) {
	logLen, _ := GlGetShaderiv(shader, INFO_LOG_LENGTH)
	buf := make([]byte, logLen)
	C.glGetShaderInfoLog(C.GLuint(shader), C.GLsizei(logLen), (*C.GLsizei)(unsafe.Pointer(&logLen)),
		(*C.GLchar)(unsafe.Pointer(&buf[0])))
	return string(buf), nil
}

func GlGetAttribLocation(prog uint, name string) (int, error) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	attrib := int(C.glGetAttribLocation(C.GLuint(prog), (*C.GLchar)(unsafe.Pointer(nameC))))
	// FIXME: check that attrib != -1
	return attrib, nil
}

func GlGetUniformLocation(prog uint, name string) (int, error) {
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))
	attrib := int(C.glGetUniformLocation(C.GLuint(prog), (*C.GLchar)(unsafe.Pointer(nameC))))
	// FIXME: check that attrib != -1
	return attrib, nil
}

func GlEnableVertexAttribArray(handle int) error {
	C.glEnableVertexAttribArray(C.GLuint(handle))
	return nil
}

func GlVertexAttribPointer(handle int, count int, glType C.GLenum, normalized bool, stride int, data uintptr) error {
	var norm C.GLboolean
	if normalized {
		norm = TRUE
	} else {
		norm = FALSE
	}
	C.glVertexAttribPointer(C.GLuint(handle), C.GLint(count), glType, norm, C.GLsizei(stride),
		(unsafe.Pointer(data)))
	return nil
}

func GlUniform2f(handle int, v1, v2 float32) error {
	C.glUniform2f(C.GLint(handle), C.GLfloat(v1), C.GLfloat(v2))
	return nil
}

func GlUniform3f(handle int, v1, v2, v3 float32) error {
	C.glUniform3f(C.GLint(handle), C.GLfloat(v1), C.GLfloat(v2), C.GLfloat(v3))
	return nil
}

func GlUniformMatrix4fv(handle int, count int, transpose bool, data []float32) error {
	var tp C.GLboolean
	if transpose {
		tp = TRUE
	} else {
		tp = FALSE
	}
	C.glUniformMatrix4fv(C.GLint(handle), C.GLsizei(count), tp,
		(*C.GLfloat)(unsafe.Pointer(&data[0])))
	return nil
}

func GlGetString(name uint) string {
	val := C.glGetString(C.GLenum(name))
	return C.GoString((*C.char)(unsafe.Pointer(val)))
}

func GlViewport(x, y, szX, szY int) error {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(szX), C.GLsizei(szY))
	return nil
}

func GlDrawArrays(mode uint, first, count int) error {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
	return nil
}
