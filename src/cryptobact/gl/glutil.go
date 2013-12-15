package gl

// utility functions around C api

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>
#include <GLES2/gl2.h>
#cgo android LDFLAGS: -lGLESv2
*/
import "C"
import "errors"
import "fmt"

func GetShaderInfoLog()  {}
func GetProgramInfoLog() {}

func LoadShader(shType uint, source string) (uint, error) {
	handle, _ := GlCreateShader(shType) // FIXME: handle error
	ErrPanic()
	GlShaderSource(handle, []string{source}) // FIXME: handle error
	ErrPanic()
	GlCompileShader(handle) // FIXME: handle error
	ErrPanic()
	isCompiled, _ := GlGetShaderiv(handle, COMPILE_STATUS) // FIXME: handle f* error!
	if isCompiled != TRUE {
		return 0, errors.New("cannot compile shader")
	}
	return handle, nil
}

func CreateProgram(vertShader, fragShader string) (uint, error) {
	vxHandle, _ := LoadShader(VERTEX_SHADER, vertShader)
	fragHandle, _ := LoadShader(FRAGMENT_SHADER, fragShader)

	prog, _ := GlCreateProgram()

	GlAttachShader(prog, vxHandle)
	ErrPanic()
	GlAttachShader(prog, fragHandle)
	ErrPanic()

	GlLinkProgram(prog)
	ErrPanic()

	linkOk, _ := GlGetProgramiv(prog, LINK_STATUS)
	ErrPanic()
	if linkOk != TRUE {
		log, _ := GlGetProgramInfoLog(prog)
	ErrPanic()
		return 0, errors.New("failed to link program: " + log)
	}
	return prog, nil
}

func MakeProjectionMatrix(left, right, bottom, top, near, far float32, matrix []float32) []float32 {
	matrix[0] = 2.0 / (right - left)
	matrix[5] = 2.0 / (top - bottom)
	matrix[10] = -2.0 / (far - near)
	matrix[12] = -(right + left) / (right - left)
	matrix[13] = -(top + bottom) / (top - bottom)
	matrix[14] = -(far + near) / (far - near)
	matrix[15] = 1.0
	return matrix
}

func ErrPanic() {
	err := C.glGetError()
	if err != C.GL_NO_ERROR {
		panic(fmt.Errorf("gl error: %v", err))
	}
}
