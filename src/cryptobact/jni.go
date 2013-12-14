package main

/*
#include <stdlib.h>
#include <jni.h>
#include <android/input.h>

*/
import "C"
import "log"
import "unsafe"

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
