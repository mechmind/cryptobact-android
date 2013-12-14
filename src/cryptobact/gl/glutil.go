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
import "unsafe"

func GetShaderInfoLog()  {}
func GetProgramInfoLog() {}

func LoadShader(shType int, source string) (uint, error) {
	handle, _ := GlCreateShader(shType)                    // FIXME: handle error
	GlShaderSource(shType, []string{source})               // FIXME: handle error
	GlCompileShader(handle)                                // FIXME: handle error
	isCompiled, _ := GlGetShaderiv(handle, COMPILE_STATUS) // FIXME: handle f* error!
	if isCompiled == 0 {
		return 0, errors.New("cannot compile shader")
	}
	return handle, nil
}

func CreateProgram(vertShader, fragShader string) (uint, error) {
	vxHandle, _ := LoadShader(VERTEX_SHADER, vertShader)
	fragHandle, _ := LoadShader(FRAGMENT_SHADER, fragHandle)
}
